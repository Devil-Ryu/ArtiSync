<template>
    <div>
        <TForm label-align="top">
            <!-- 文章设置卡片 -->
            <SettingCard title="文章设置">
                <template #content>
                    <t-row style="width: 100%;" :gutter="10">
                        <t-form-item label="图片读取方式" name="imageRead" style="width: 100%;">
                            <t-col flex="none">
                                <TSelect v-model:value="configStore.systemConfig.imageReadType"
                                    :options="imageReadOptions" />
                            </t-col>
                        </t-form-item>
                    </t-row>
                    <t-row style="width: 100%;" :gutter="10">
                        <t-form-item label="图片文件夹 ( 本地文件读取时生效 )" name="imageDir" style="width: 100%;">
                            <t-col flex="none">
                                <TSelect v-model:value="configStore.systemConfig.imageSelect"
                                    :options="imagePathOptions" />
                            </t-col>
                            <t-col flex="auto">
                                <t-input v-model="configStore.systemConfig.imagePath" placeholder="相对文章时为空则默认跟文章目录同一级">
                                    <template #suffixIcon>
                                        <FolderOpen1Icon style="cursor: pointer" @click="SetDir('imagePath')" />
                                    </template>
                                </t-input>
                            </t-col>
                        </t-form-item>
                    </t-row>


                </template>
            </SettingCard>

            <!-- 系统设置 -->
            <SettingCard title="系统设置">
                <template #content>
                    <t-row :gutter="10">
                        <t-col flex="auto">
                            <t-form-item label="系统配置路径" name="sysPath">
                                <t-input v-model="configStore.sysPath" placeholder="请输入" disabled>
                                    <!-- <template #suffixIcon>
                                    <FolderOpen1Icon style="cursor: pointer" @click="SetDir('sysPath')" />
                                </template> -->
                                </t-input>
                            </t-form-item>
                        </t-col>
                    </t-row>
                    <t-row :gutter="10">
                        <t-col flex="auto">
                            <t-form-item label="数据库文件" name="dbPath">
                                <t-input v-model="configStore.systemConfig.dbPath" placeholder="请输入">
                                    <template #suffixIcon>
                                        <FolderOpen1Icon style="cursor: pointer" @click="SetDir('dbPath')" />
                                    </template>
                                </t-input>
                            </t-form-item>
                        </t-col>
                    </t-row>
                    <t-row :gutter="10">
                        <t-col flex="auto">
                            <t-form-item label="代理" name="proxy">
                                <t-input v-model="configStore.systemConfig.proxyURL"
                                    placeholder="http://127.0.0.1:8080">
                                </t-input>
                            </t-form-item>
                        </t-col>
                    </t-row>
                    <t-row :gutter="10">
                        <t-col flex="auto">
                            <t-form-item label="接口请求速度限制" name="requestSleep">
                                <t-input-number style="width: 100%;" v-model="configStore.systemConfig.requestSleep"
                                    align="right" label="每运行一个接口睡眠" suffix="秒" theme="column" default-value="1"
                                    placeholder="1" />
                            </t-form-item>
                        </t-col>
                    </t-row>
                </template>
            </SettingCard>

            <!-- 保存按钮 -->
            <t-space style="display: flex;justify-content: center">
                <t-button variant="outline" theme="danger" @click="resetSystemConfig">
                    <template #icon>
                        <RefreshIcon />
                    </template>
                    重置
                </t-button>
                <t-button variant="outline" theme="primary" @click="saveSystemConfig">
                    <template #icon>
                        <SaveIcon />
                    </template>
                    保存
                </t-button>
            </t-space>
        </TForm>
    </div>

</template>

<script setup>
import { MessagePlugin } from "tdesign-vue-next";
import { Connect, GetConfigFilePath, LoadJSONFile, SaveJSONFile } from "@/wailsjs/go/api/DBController.js";
import SettingCard from "./components/SettingCard.vue"
import { onMounted, ref } from "vue";
import { InitConfig, OpenDir } from "@/wailsjs/go/api/ATController";
import { FolderOpen1Icon, RefreshIcon, SaveIcon } from "tdesign-icons-vue-next";
import { useConfigStore } from "@/src/store/config"
import CodeViewer from "../Components/CodeViewer.vue";

const configStore = useConfigStore()

const imagePathOptions = [
    { label: "相对文章目录", value: "相对文章目录" },
    { label: "固定图片目录", value: "固定图片目录" },
]

const imageReadOptions = [
    { label: "不读取", value: "NONE" },
    { label: "本地文件读取", value: "LOCAL" },
    { label: "远程链接下载", value: "REMOTE" },
]

function saveSystemConfig() {
    GetConfigFilePath().then((configFilePath) => {
        SaveJSONFile(configFilePath, configStore.systemConfig).then(() => {
            InitConfig().then(() => {
                MessagePlugin.success("保存成功")
            }).catch((erri) => {
                MessagePlugin.error("保存失败: " + erri)
            })
        }).catch((err) => {
            MessagePlugin.error("保存失败: " + err)
        })
    }).catch((err) => {
        MessagePlugin.error("获取配置文件失败: " + err)
    })
}

function resetSystemConfig() {
    GetConfigFilePath().then((configFilePath) => {
        SaveJSONFile(configFilePath, configStore.defaultConfig).then(() => {
            InitConfig().then(() => {
                configStore.systemConfig = JSON.parse(JSON.stringify(configStore.defaultConfig))  // 需要拷贝
                MessagePlugin.success("重置成功")
            }).catch((erri) => {
                MessagePlugin.error("重置失败: " + erri)
            })
        }).catch((err) => {
            MessagePlugin.error("重置失败: " + err)
        })
    }).catch((err) => {
        MessagePlugin.error("获取配置文件失败: " + err)
    })

}

function SetDir(valueName) {
    OpenDir(configStore.systemConfig[valueName]).then((response) => {
        if (response.StatusCode == 200 && response.Data !== "") {
            var selectedDir = response.Data
            configStore.systemConfig[valueName] = selectedDir
        }
    }).catch((err)=>{
        MessagePlugin.error("打开文件夹错误: "+err)
    })
}

</script>

<style scoped></style>