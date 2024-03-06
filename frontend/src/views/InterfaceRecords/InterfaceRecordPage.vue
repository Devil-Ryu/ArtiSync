<template>
  <div style="height: 100%; " ref="viewRef">
    <t-form>
      <t-row>
        <t-col :span="12">
          <t-row :gutter="[10, 24]">
            <t-col :span="4">
              <t-form-item label="请求时间" name="datetime">
                <t-date-range-picker v-model="filterValue.date_time" :presets="presets" enable-time-picker />
              </t-form-item>
            </t-col>
            <t-col :span="4">
              <t-form-item label="接口名称" name="interface">
                <t-input type="search" placeholder="请输入接口名称" v-model="filterValue.name" />
              </t-form-item>
            </t-col>
            <t-col :span="4">
              <t-form-item label="运行状态" name="status">
                <t-select multiple :min-collapsed-num="1" :options="statusOptions" v-model="filterValue.status" />
              </t-form-item>
            </t-col>
          </t-row>
          <t-row :gutter="[10, 24]" style="margin-top: 24px;">
            <t-col :span="4">
              <t-form-item type="search" placeholder="请输入平台名称" label="平台名称" name="platforName">
                <t-input v-model="filterValue.platform_name" />
              </t-form-item>
            </t-col>

            <t-col :span="4">
              <t-form-item label="记录ID" name="recordID">
                <t-input type="search" placeholder="请输入记录ID" v-model="filterValue.record_id" />
              </t-form-item>
            </t-col>
            <t-col :span="3">
              <t-space :size="10">
                <t-button variant="outline" @click="columnControllerVisible = true">列配置</t-button>
                <t-button type="reset" variant="outline" theme="warning" @click="reset"> 重置 </t-button>
                <t-button type="reset" variant="outline" theme="primary" @click="query"> 查询 </t-button>
                <t-popconfirm theme="danger" content="操作将删除所有记录，确认删除吗?" @confirm="clear">
                  <t-button type="reset" variant="outline" theme="danger"> 清空 </t-button>
                </t-popconfirm>
              </t-space>
            </t-col>
          </t-row>
        </t-col>
      </t-row>
    </t-form>
    <t-table :max-height="tableHeight" :columns="columns" :data="data" :pagination="pagination" row-key="ID" lazy-load
      v-model:displayColumns="displayColumns" :column-controller="columnControllerConfig"
      v-model:columnControllerVisible="columnControllerVisible" resizable style="margin-top: 10px; "
      @page-change="onPageChange">
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
        <t-link theme="danger" @click="deletRecord(row)">删除</t-link>
      </template>
    </t-table>
    <RecordDetailDialog />
  </div>
</template>

<script setup>
import { computed, onActivated, onMounted, onUnmounted, ref, watch } from "vue";
import { DeleteInterfaceRecord, QueryInterfaceRecords, ClearInterfaceRecord, GetInterfaceRecords } from "@/wailsjs/go/api/DBController";
import { statusNameListMap } from "@/src/store/article"
import { useInterfaceRecordsStore } from "@/src/store/platform";
import RecordDetailDialog from "@/src/views/Components/RecordDetailDialog/RecordDetailDialog.vue";
import { MessagePlugin } from "tdesign-vue-next";
import { EventsOn } from "@/wailsjs/runtime/runtime";
import dayjs from "dayjs"
const viewRef = ref()
const tableHeight = ref(490)
const interfaceRecordsStore = useInterfaceRecordsStore()
const isLoading = ref(false)
const filterValue = ref({
  date_time: [],
});
// 根据当前页面显示过滤器
// const filterValue = computed(() => {
//   return interfaceRecordsStore.filters[interfaceRecordsStore.curFilter]
// })

// 接口运行状态部分参数（TODO（增加列控制，增加清空功能））

// 运行状态选项
const statusOptions = [
  { value: "运行成功", label: "运行成功" },
  { value: "运行失败", label: "运行失败" },
  { value: "运行中", label: "运行中" },
]

// 日期选择框标签
const presets = ref({
  最近7天: [dayjs().hour(0).minute(0).second(0).subtract(6, 'day'), dayjs().hour(23).minute(59).second(59)],
  最近3天: [dayjs().hour(0).minute(0).second(0).subtract(2, 'day'), dayjs().hour(23).minute(59).second(59)],
  今天: [dayjs().hour(0).minute(0).second(0), dayjs().hour(23).minute(59).second(59)],
});

const columns = ref([
  { colKey: 'DateTime', title: '请求时间', width: 180, ellipsis: true },
  { colKey: 'RecordID', title: '记录ID', width: 100, ellipsis: true },
  { colKey: 'ArticleName', title: '文章名称', width: 100, ellipsis: true },
  { colKey: 'PlatformName', title: '平台名称', width: 100, ellipsis: true },
  { colKey: 'Tag', title: '标签', width: 70, align: 'center', ellipsis: true, },
  { colKey: 'Name', title: '接口名称', width: 100, ellipsis: true, },
  { colKey: 'RequestURL', title: '接口URL', ellipsis: true, },
  // { colKey: 'ResponseMessage', title: '响应内容', ellipsis: true, },
  { colKey: 'Status', title: '运行状态', width: 100, align: 'center', ellipsis: true, },
  { colKey: 'operation', title: '操作', width: 80, fixed: "right", align: 'center' },
])

const columnControllerVisible = ref(false);
const columnControllerConfig = computed(() => ({
  // 隐藏组件内部的 列配置按钮
  hideTriggerButton: true,
  // 用于设置允许用户对哪些列进行显示或隐藏的控制，默认为全部字段
  fields: ['DateTime', 'RecordID','ArticleName','PlatformName' , 'Tag', 'Name', 'RequestURL', 'Status'],

  // 弹框组件属性透传
  dialogProps: { preventScrollThrough: true },
  // 列配置按钮组件属性透传
  // buttonProps: customText.value ? { content: '显示列控制', theme: 'primary', variant: 'base' } : undefined,
  buttonProps: undefined,


}))

const staticColumn = ['Name', 'operation'];
const displayColumns = ref(staticColumn.concat(['DateTime', 'RecordID', 'RequestURL', 'Tag', 'Status']));

const data = ref([])
const pagination = ref({
  defaultCurrent: 1,
  defaultPageSize: 10,
  total: 0,
  current: 1,
})

function deletRecord(row) {
  DeleteInterfaceRecord(row).then(result => {
    let index = data.value.findIndex((item) => item.ID == row.ID)
    data.value.splice(index, 1)
    pagination.value.total -= 1
    MessagePlugin.success("删除成功")
  }).catch((err) => {
    MessagePlugin.error("删除失败: " + err)
  })
}

function openRecordDetail(row) {
  // 根据ID查询单个记录
  GetInterfaceRecords({ "id": row.ID }).then((result) => {
    interfaceRecordsStore.curRecord = result[0]
    interfaceRecordsStore.detailDialogVisible = true
  }).catch(err => {
    MessagePlugin.err("获取记录失败: " + err)
  })

}

function reset() {
  filterValue.value = { date_time: [] }
  fetchData({ current: pagination.value.defaultCurrent, pageSize: pagination.value.defaultPageSize })
}

function clear() {
  ClearInterfaceRecord().then(() => {
    fetchData({ current: pagination.value.defaultCurrent, pageSize: pagination.value.defaultPageSize })
    MessagePlugin.success("已删除全部记录")
  }).catch((err) => {
    MessagePlugin.error("清空错误：" + err)
  })
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
      pagination.value.pageSize = result.pageSize;
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


// 函数区
EventsOn("UpdateInterfaceRecord", (recordID) => {
  // 如果是文章详情页面
  query()
})


// TODO（解决表格百分比高度问题）
onMounted(() => {
  console.log("onMounted: 接口记录页面")
  filterValue.value = interfaceRecordsStore.filters
  query()
})

const unsubscribe = interfaceRecordsStore.$onAction(({ name, store, args, after, onError }) => {
  console.log("接口刷新调用: ", name)
  after((filters) => {
    console.log("刷新过滤器: ", filters)
    filterValue.value = interfaceRecordsStore.filters
    query()
  })

  // 如果 action 抛出或返回一个拒绝的 promise，这将触发
  onError((error) => {
      console.warn(
        `接口刷新调用失败 "${name}" after ${Date.now() - startTime}ms.\nError: ${error}.`
      )
    })

})

</script>

<style></style>