<template>
<t-table :columns="columns" :data="interfaceRecordStore.testCaches" rowKey="Serial">
    <template #operation="{ row }">
        <t-space :size="5">
            <t-link theme="primary" @click="openDialog(row)">详情</t-link>
            <t-link theme="danger" @click="deleteCache(row)">删除</t-link>
        </t-space>
      </template>
</t-table>
<t-dialog v-model:visible="dialogVisible" width="80%" header="缓存值" cancelBtn="关闭" :confirmBtn="null">
    <div style="overflow: auto; max-height: 400px;">
        {{ dialogContent }}
    </div>
</t-dialog>
</template>
<script setup>
import {ref} from "vue"
import { EventsOn } from "@/wailsjs/runtime/runtime";
import { useInterfaceRecordsStore } from "@/src/store/platform"
import { DeleteTestNetControllerCache} from "@/wailsjs/go/api/ATController";
import { MessagePlugin } from "tdesign-vue-next";
const interfaceRecordStore = useInterfaceRecordsStore()

const dialogVisible = ref(false)
const dialogContent = ref("")

const columns = ref([
  { colKey: 'Serial', title: '接口编号', width: 120, ellipsis: true, },
  { colKey: 'Name', title: '接口名称', width: 100, ellipsis: true, },
  { colKey: 'Reponse', title: '响应内容', ellipsis: true, },
  { colKey: 'ReponseType', title: '响应类型', width: 100, align: 'center', ellipsis: true, },
  { colKey: 'operation', title: '操作', width: 100, fixed: "right", align: 'center'},
])

function openDialog(row) {
    dialogContent.value = row.Reponse
    dialogVisible.value = true
}

function deleteCache(row) {
    DeleteTestNetControllerCache(row.Serial).then(()=>{
        let index = interfaceRecordStore.testCaches.findIndex((item)=> item.Serial === row.Serial)
        interfaceRecordStore.testCaches.splice(index, 1)
        MessagePlugin.success("删除成功")

    }).catch((err)=>{
        MessagePlugin.error("删除失败: "+err)
    })
}

EventsOn("UpdateTestNetworkPool", (result) => {
    console.log("UpdateTestNetworkPool", result)
  interfaceRecordStore.testCaches = result
})
</script>

<style scoped>

</style>1