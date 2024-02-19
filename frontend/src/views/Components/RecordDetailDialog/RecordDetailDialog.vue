<template>
    <t-dialog 
        header="记录详情"
        v-model:visible="interfaceRecordsStore.detailDialogVisible"
        :confirm-btn="null"
        cancel-btn="关闭"
        @close="onColose"
        width="80%"
    >
    <template #header>
        <t-row class="detailHeader">
            <t-col class="title"><t-tag :theme="statusNameListMap[interfaceRecordsStore.curRecord.Status].theme">{{ statusNameListMap[interfaceRecordsStore.curRecord.Status].label }}</t-tag></t-col>
            <t-col class="content">{{ interfaceRecordsStore.curRecord.RequestURL }}</t-col>
        </t-row>
    </template>
            <div >
                <!-- <t-row class="detailHeader">
                    <t-col class="title"><t-tag :theme="statusThemeListMap[logDetailStore.logData.level]">{{ statusNameListMap[logDetailStore.logData.level] }}</t-tag></t-col>
                    <t-col class="content">{{logDetailStore.logData.tag}}</t-col>
                </t-row> -->
                <t-row class="detailRow">
                    <t-col class="title">接口编号: </t-col>
                    <t-col class="content">{{  interfaceRecordsStore.curRecord.Serial }}</t-col>
                </t-row>
                <t-row class="detailRow">
                    <t-col class="title">接口名称: </t-col>
                    <t-col class="content">{{  interfaceRecordsStore.curRecord.Name }}</t-col>
                </t-row>
                <t-row class="detailRow">
                    <t-col class="title">请求时间: </t-col>
                    <t-col class="content">{{ interfaceRecordsStore.curRecord.DateTime }}</t-col>
                </t-row>
                <t-row class="detailRow">
                    <t-col class="title">触发模块: </t-col>
                    <t-col class="content">[{{ interfaceRecordsStore.curRecord.ArticleName }}] - [{{ interfaceRecordsStore.curRecord.PlatformName }}]</t-col>
                </t-row>
                <t-row class="detailRow">
                    <t-col class="title">记录ID: </t-col>
                    <t-col class="content">{{  interfaceRecordsStore.curRecord.RecordID }}</t-col>
                </t-row>
                <t-row justify="space-between">
                    <t-col>
                        <t-radio-group v-model:value="bodyType">
                            <t-radio-button value="Request" label="请求报文"></t-radio-button>
                            <t-radio-button value="Response"  label="响应报文"></t-radio-button>
                        </t-radio-group>
                    </t-col>
                    <t-col>
                        <t-select :options="['UTF-8', 'GBK'].map((value) => ({ label: value, value }))" v-model="contentType" @change="contentTypeChange"></t-select>
                    </t-col>
                </t-row>
                <div style="background-color: rgb(247, 248, 250); border-radius: 10px; margin-top: 20px; padding: 10px;">
                    <net-message-viewer :content="bodyContent" :content-type="contentType"/>
                </div>
            </div>
        <!-- </t-card> -->
    </t-dialog>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { useInterfaceRecordsStore } from "@/src/store/platform";
import {statusNameListMap} from "@/src/store/article"
import NetMessageViewer from "./components/NetMessageViewer.vue"

const interfaceRecordsStore = useInterfaceRecordsStore()
const contentType  = ref('UTF-8')
const bodyType = ref('Request')
const bodyOptions = [
    { value: "Request", label: "请求报文"},
    { value: "Response", label: "响应报文"},
]
const bodyContent = computed(()=>{
    if (bodyType.value == "Request") {
        console.log("Request")
        return interfaceRecordsStore.curRecord.RequestMessage
    } else {
        console.log("Response")
        return interfaceRecordsStore.curRecord.ResponseMessage
    }
})


const emits = defineEmits(['close'])

function onColose() {
    contentType.value = "UTF-8"
    bodyType.value = "Request"
    emits('close')
}

function contentTypeChange(value) {
    // 尝试JSON解码
    if (value === "JSON") {
        try {
            JSON.parse(logDetailStore.logData.message)
        } catch (error) {
            contentType.value = "Request"
            MessagePlugin.error("JSON解析失败: "+error)
        }
    }
}

</script>

<style scoped>
.detailHeader {
    display: flex;
    width: 100%;
    align-items: center;
    justify-items: center;
    font-size: 0.8125rem;

    .content {
        margin-left: 12px;
        font-size: 14px;
        font-weight: 500;
        word-break: "break-all";
    }
}

.detailRow {
    display: flex;
    width: 100%;
    align-items: center;

    .title {
        width: 80px;
        color: rgba(0, 0, 0, 0.5);
        font-size: 14px;
        flex-shrink: 0;
        max-width: 80px ;
    }

    .content {
        font-size: 14px;
        font-weight: 500;
        word-break: "break-all";
        /* margin-right: 16px; */
    }
    
}

.detailContent {
    margin-top: 5px;
}



</style>