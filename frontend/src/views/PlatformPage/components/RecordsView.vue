<template>
    <t-table :columns="columns" :data="data"  :pagination="pagination" @page-change="onPageChange" row-key="ID" lazy-load style="margin-top: 10px; ">
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
            <t-link theme="primary">详情</t-link>
            <t-link theme="danger">删除</t-link>
        </template>
    </t-table>
</template>

<script setup>
import { ref } from "vue"
import { statusNameListMap } from "@/src/store/article"
import { useInterfacesStore } from "@/src/store/platform";

const interfaceStore = useInterfacesStore()

const data = ref([])
const columns = ref([
    { colKey: 'DateTime', title: '请求时间', width: 180, ellipsis: true },
    { colKey: 'RecordID', title: '记录ID', width: 100, ellipsis: true },
    { colKey: 'Tag', title: '标签', width: 70, align: 'center', ellipsis: true, },
    { colKey: 'Name', title: '接口名称', width: 100, ellipsis: true, },
    { colKey: 'RequestURL', title: '接口URL', ellipsis: true, },
    { colKey: 'ResponseMessage', title: '响应内容', ellipsis: true, },
    { colKey: 'Status', title: '运行状态', width: 100, align: 'center', ellipsis: true, },
    { colKey: 'operation', title: '操作', width: 80, fixed: "right", align: 'center' },
])

const pagination = ref({
  defaultCurrent: 1,
  defaultPageSize: 10,
  total: 0,
  current: 1,
})

function fetchData(pageInfo) {
  try {
    isLoading.value = true;
    const { current, pageSize } = pageInfo;
    QueryInterfaceRecords({"recod_id":"TEST", "platform_name":interfaceStore.platform.Name}, current, pageSize).then((result) => {
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

</script>

<style scoped></style>