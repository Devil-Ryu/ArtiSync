<template>
    <t-dialog 
        header="日志详情"
        v-model:visible="logDetailStore.visible"
        :confirm-btn="null"
        cancel-btn="关闭"
        @close="onColose"
        width="80%"
    >
    <template #header>
        <t-row class="detailHeader">
            <t-col class="title"><t-tag :theme="statusThemeListMap[logDetailStore.logData.level]">{{ statusNameListMap[logDetailStore.logData.level] }}</t-tag></t-col>
            <t-col class="content">{{logDetailStore.logData.tag}}</t-col>
        </t-row>
    </template>
        <!-- <t-card :shadow="true"> -->
            <div >
                <!-- <t-row class="detailHeader">
                    <t-col class="title"><t-tag :theme="statusThemeListMap[logDetailStore.logData.level]">{{ statusNameListMap[logDetailStore.logData.level] }}</t-tag></t-col>
                    <t-col class="content">{{logDetailStore.logData.tag}}</t-col>
                </t-row> -->
                <t-row class="detailRow">
                    <t-col class="title">日志级别: </t-col>
                    <t-col class="content">{{  statusNameListMap[logDetailStore.logData.level] }}</t-col>
                </t-row>
                <t-row class="detailRow">
                    <t-col class="title">日志标签: </t-col>
                    <t-col class="content">{{logDetailStore.logData.tag}}</t-col>
                </t-row>
                <t-row class="detailRow">
                    <t-col class="title">日志时间: </t-col>
                    <t-col class="content">{{ logDetailStore.logData.datetime }}</t-col>
                </t-row>
                <t-row justify="space-between" style="margin-top: 10px;">
                    <t-col>
                        <t-radio-group default-value="日志详情">
                            <!-- <t-radio-button value="请求报文" label="请求报文" ></t-radio-button> -->
                            <t-radio-button value="日志详情"  label="日志详情"></t-radio-button>
                        </t-radio-group>
                    </t-col>
                    <t-col>
                        <t-select :options="['TEXT', 'JSON'].map((value) => ({ label: value, value }))" v-model="contentType" @change="contentTypeChange"></t-select>
                    </t-col>
                </t-row>
                
                <div style="background-color: rgb(247, 248, 250); border-radius: 10px; margin-top: 20px;">
                    <!-- {{ logDetailStore.logData.message }} -->
                    <code-viewer :content="logDetailStore.logData.message" :content-type="contentType"/>
                </div>
            </div>
        <!-- </t-card> -->
    </t-dialog>
</template>

<script setup>
import {useLogDetailStore, statusNameListMap, statusThemeListMap} from "@/src/store/log"
import { onMounted, ref } from "vue";
import CodeViewer from "../../Components/CodeViewer.vue";
import { MessagePlugin } from "tdesign-vue-next";
const logDetailStore = useLogDetailStore()
const contentType  = ref('TEXT')
const emits = defineEmits(['close'])

function onColose() {
    contentType.value = "TEXT"
    emits('close')
}

function contentTypeChange(value) {
    // 尝试JSON解码
    if (value === "JSON") {
        try {
            JSON.parse(logDetailStore.logData.message)
        } catch (error) {
            contentType.value = "TEXT"
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
    margin-top: 10px;

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