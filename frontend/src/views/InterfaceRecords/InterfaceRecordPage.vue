<template>
  <div>
    <t-form>
      <t-row>
        <t-col :span="12">
          <t-row :gutter="[10, 24]" >
            <t-col :span="6">
              <t-form-item label="请求时间" name="datetime">
                <t-date-range-picker v-model="filterValue.date_time" enable-time-picker />
              </t-form-item>
            </t-col>
            <t-col :span="6">
              <t-space :size="10" >
                <t-button type="reset" variant="base" theme="default" @click="query"> 查询 </t-button>
                <t-button type="reset" variant="base" theme="default" @click="reset"> 重置 </t-button>
                <t-button type="reset" variant="base" theme="default" @click="reset"> 清空 </t-button>
              </t-space>
            </t-col>
          </t-row>
          <t-row :gutter="[10, 24]" style="margin-top: 24px;">
            <t-col :span="3">
              <t-form-item label="平台名称" name="platforName">
                <t-select multiple :min-collapsed-num="1" :options="platformNameOptions"
                  v-model="filterValue.platform_name" />
              </t-form-item>
            </t-col>
            <t-col :span="3">
              <t-form-item label="接口名称" name="interface">
                <t-input type="search" placeholder="请输入接口名称" v-model="filterValue.name" />
              </t-form-item>
            </t-col>
            <t-col :span="3">
              <t-form-item label="运行状态" name="status">
                <t-select multiple :min-collapsed-num="1" :options="statusOptions" v-model="filterValue.status" />
              </t-form-item>
            </t-col>
            <t-col :span="3">
              <t-form-item label="记录ID" name="recordID">
                <t-input type="search" placeholder="请输入记录ID" v-model="filterValue.record_id" />
              </t-form-item>
            </t-col>
          </t-row>
        </t-col>
        <!-- <t-col  style="margin-left: 20px;">
          <t-button type="reset" variant="base" theme="default" @click="query"> 查询 </t-button>
          <t-button type="reset" variant="base" theme="default" @click="reset"> 重置 </t-button>
        </t-col> -->
      </t-row>
    </t-form>
    <t-table :columns="columns" :data="data" :pagination="pagination" row-key="ID" lazy-load resizable
      style="margin-top: 10px;" @page-change="onPageChange">
      <template #Status="{ row }">
        <t-tag variant="light-outline" :theme="statusNameListMap[row.Status].theme">
          {{ statusNameListMap[row.Status].label }}
        </t-tag>
      </template>
      <template #Tag="{ row }">
        <t-tag variant="light-outline" theme="primary">
          {{ row.Tag }}
        </t-tag>
      </template>
      <template #operation="{ row }">
        <t-link theme="primary" @click="openRecordDetail(row)">详情</t-link>
      </template>
    </t-table>
    <RecordDetailDialog />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { GetInterfaceRecords, QueryInterfaceRecords } from "@/wailsjs/go/api/DBController";
import { statusNameListMap } from "@/src/store/article"
import { useInterfaceRecordsStore } from "@/src/store/platform";
import RecordDetailDialog from "@/src/views/Components/RecordDetailDialog/RecordDetailDialog.vue";

const interfaceRecordsStore = useInterfaceRecordsStore()
const isLoading = ref(false)
const filterValue = ref({
  date_time: [],
});

// 接口运行状态部分参数（TODO（增加列控制，增加筛选功能，增加删除功能，增加清空功能））
// 平台名称选项
const platformNameOptions = computed(() => {
  var tmp = []
  data.value.forEach(item => {
    tmp.push(item.PlatformName)
  })
  var s = new Set(tmp)
  var result = []
  s.forEach(key => {
    result.push({ value: key, label: key })
  })
  return result
})
// 运行状态选项
const statusOptions = [
  { value: "运行成功", label: "运行成功" },
  { value: "运行失败", label: "运行失败" },
  { value: "运行中", label: "运行中" },
]

// 日期选择框标签
// const presets = ref({
//   最近7天: [dayjs().subtract(6, 'day'), dayjs()],
//   最近3天: [dayjs().subtract(2, 'day'), dayjs()],
//   今天: [dayjs(), dayjs()],
// });

const columns = ref([
  { colKey: 'DateTime', title: '请求时间', width: 180, ellipsis: true },
  { colKey: 'RecordID', title: '记录ID', width: 120, ellipsis: true },
  { colKey: 'Tag', title: '标签', width: 80, align: 'center', ellipsis: true, },
  { colKey: 'Name', title: '接口名称', width: 100, ellipsis: true, },
  { colKey: 'RequestURL', title: '接口URL', ellipsis: true, },
  { colKey: 'ResponseMessage', title: '响应内容', ellipsis: true, },
  { colKey: 'Status', title: '运行状态', width: 100, align: 'center', ellipsis: true, },
  { colKey: 'operation', title: '操作', width: 80, fixed: "right", align: 'center' },
])

const data = ref([])
const pagination = ref({
  defaultCurrent: 1,
  defaultPageSize: 10,
  total: 100,
  current: 1,
})



function openRecordDetail(row) {
  interfaceRecordsStore.curRecord = row
  interfaceRecordsStore.detailDialogVisible = true
}

function reset() {
  filterValue.value = {}
  fetchData({ current: pagination.value.defaultCurrent, pageSize: pagination.value.defaultPageSize })
}

function query() {
  fetchData({ current: pagination.value.defaultCurrent, pageSize: pagination.value.defaultPageSize })
}
function fetchData(pageInfo) {
  try {
    isLoading.value = true;
    const { current, pageSize } = pageInfo;
    QueryInterfaceRecords(filterValue.value, current, pageSize).then((result) => {
      data.value = result.result
      pagination.value.current = result.pageNum;
      pagination.value.total = result.totalRows;
    })
    isLoading.value = false;
  } catch (err) {
    console.log("in err", err);
    data.value = [];
  }
};

function onPageChange(pageInfo) {
  fetchData(pageInfo)
}

onMounted(() => {
  fetchData({ current: pagination.value.defaultCurrent, pageSize: pagination.value.defaultPageSize })
})
</script>

<style></style>