<template>
    <t-dialog header="记录详情" v-model:visible="interfaceRecordsStore.detailDialogVisible" :confirm-btn="null" cancel-btn="关闭"
        @close="onColose" width="80%">
        <template #header>
            <t-row class="detailHeader">
                <t-col class="title"><t-tag :theme="statusNameListMap[interfaceRecordsStore.curRecord.Status].theme">{{
                    statusNameListMap[interfaceRecordsStore.curRecord.Status].label }}</t-tag></t-col>
                <t-col class="content">{{ interfaceRecordsStore.curRecord.RequestURL }}</t-col>
            </t-row>
        </template>
        <div>
            <!-- <t-row class="detailHeader">
                    <t-col class="title"><t-tag :theme="statusThemeListMap[logDetailStore.logData.level]">{{ statusNameListMap[logDetailStore.logData.level] }}</t-tag></t-col>
                    <t-col class="content">{{logDetailStore.logData.tag}}</t-col>
                </t-row> -->
            <t-row class="detailRow">
                <t-col class="title">接口编号: </t-col>
                <t-col class="content">{{ interfaceRecordsStore.curRecord.Serial }}</t-col>
            </t-row>
            <t-row class="detailRow">
                <t-col class="title">接口名称: </t-col>
                <t-col class="content">{{ interfaceRecordsStore.curRecord.Name }}</t-col>
            </t-row>
            <t-row class="detailRow">
                <t-col class="title">请求时间: </t-col>
                <t-col class="content">{{ interfaceRecordsStore.curRecord.DateTime }}</t-col>
            </t-row>
            <t-row class="detailRow">
                <t-col class="title">触发模块: </t-col>
                <t-col class="content">[{{ interfaceRecordsStore.curRecord.ArticleName }}] - [{{
                    interfaceRecordsStore.curRecord.PlatformName }}]</t-col>
            </t-row>
            <t-row class="detailRow">
                <t-col class="title">记录ID: </t-col>
                <t-col class="content">{{ interfaceRecordsStore.curRecord.RecordID }}</t-col>
            </t-row>
            <t-row justify="space-between">
                <t-col>
                    <t-radio-group v-model:value="bodyType">
                        <t-radio-button value="Request" label="请求报文"></t-radio-button>
                        <t-radio-button value="Response" label="响应报文"></t-radio-button>
                    </t-radio-group>

                </t-col>
                <t-col>
                    <div style="display: flex;align-items: center;">
                        <t-tooltip content="复制报文">
                            <FileCopyIcon size="20" style="margin-right: 10px; color: gray; cursor: pointer;"
                            @click="copyToClipboard" />
                        </t-tooltip>
                       
                        <t-select :options="['UTF-8', 'GBK'].map((value) => ({ label: value, value }))"
                            v-model="contentType" @change="contentTypeChange"></t-select>
                    </div>
                </t-col>
            </t-row>
            <div style="background-color: rgb(247, 248, 250); border-radius: 10px; margin-top: 20px; padding: 10px;">
                <net-message-viewer :content="bodyContent" :content-type="contentType" />
            </div>
        </div>
        <!-- </t-card> -->
    </t-dialog>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { Message, MessagePlugin } from "tdesign-vue-next";
import { useInterfaceRecordsStore } from "@/src/store/platform";
import { statusNameListMap } from "@/src/store/article"
import NetMessageViewer from "./components/NetMessageViewer.vue"
import { FileCopyIcon } from "tdesign-icons-vue-next";

const interfaceRecordsStore = useInterfaceRecordsStore()
const contentType = ref('UTF-8')
const bodyType = ref('Request')
const bodyContent = computed(() => {
    if (bodyType.value == "Request") {
        return interfaceRecordsStore.curRecord.RequestMessage
    } else {
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
            MessagePlugin.error("JSON解析失败: " + error)
        }
    }
}

const authentication = () => {
    if ("clipboard" in navigator) {
        return new Promise((resolve, reject) => {
            navigator.permissions.query({ name: "clipboard-read" }).then(
                (result) => {
                    if (result.state == "granted" || result.state == "prompt") {
                        resolve(true);
                    } else {
                        resolve(false);
                    }
                },
                (error) => {
                    reject(error);
                }
            );
        });
    } else {
        alert("该浏览器暂不支持，请使用最新版本的GoogleChrome浏览");
        return Promise.resolve(false);
    }
};

function copyToClipboard() {
    navigator.clipboard.writeText(bodyContent.value).then(
        () => {
            MessagePlugin.success("复制成功")
        },
        (err) => {
            MessagePlugin.error("复制失败"+err)
        }
    );

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
        max-width: 80px;
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