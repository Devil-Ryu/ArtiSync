<template>
  <div>
    <div style="position: relative">
      <div v-for="(platform, index) in platformStore.platforms">
        <t-card :shadow="true" class="box-card">
          <template #header>
            <t-space :size="8">
              {{ platform.Name }}
              <t-tag :theme=statusMap[platform.Disabled].color variant="light">{{ statusMap[platform.Disabled].label }}
              </t-tag>
            </t-space>
            <t-space :size="8">
              <t-button variant="outline" @click="exportPlatform(platform)">导出</t-button>
              <t-button variant="outline" @click="editPlatform(platform)">编辑</t-button>
              <t-button variant="outline" @click="updatePlatform(platform)">{{ statusMap[!platform.Disabled].label }}</t-button>
              <t-button theme="danger" variant="base" @click="deletePlatform(index, platform)">删除</t-button>
            </t-space>
          </template>
        </t-card>
      </div>
      <t-space style="display: flex;justify-content: center">
        <!-- <t-button  variant="outline" theme="success" >
          <template #icon><RocketIcon /></template>
          测速
        </t-button> -->
        <t-button variant="outline" theme="primary" @click="addPlatform">
          <template #icon><AddIcon /></template>
          新增
        </t-button>
        <t-button variant="outline" theme="primary" @click="importPlatform">
          <template #icon><AddIcon /></template>
          导入
        </t-button>
      </t-space>
    </div>
    <PlatformConfigDialog />
  </div>
  

</template>
<script setup>
import {onActivated, onMounted, ref} from "vue";
import {AddIcon, RocketIcon} from "tdesign-icons-vue-next";
import {CreatePlatform, DeletePlatforms, GetPlatform, GetPlatforms, UpdatePlatform, ExportPlatform, ImportPlatform} from "../../../wailsjs/go/api/DBController.js";
import { OpenFile, SaveFile } from "@/wailsjs/go/api/ATController.js";
import PlatformConfigDialog from "./PlatformConfigDialog.vue";
import { MessagePlugin } from "tdesign-vue-next";
import { usePlatformStore, useInterfacesStore } from "@/src/store/platform"

const platformStore = usePlatformStore()
const interfacesStore = useInterfacesStore()

// const platforms = ref([])
// let platformConfig = ref({
//   visible: false,
//   title: "平台配置",
//   configData: {
//     Name: "",
//     ID: undefined,
//     Interfaces: []
//   }
// })
const statusMap = {
  true: {label: "禁用", color: "default"},
  false: {label: "启用", color: "success"}
}
// 加载平台
function loadPlatforms() {
  GetPlatforms().then(result => {
    console.log("GetPlatforms->", result)
    platformStore.platforms = result
  })
}

// 新增平台
function addPlatform() {
  CreatePlatform([{name: '新平台'}]).then(result => {
    platformStore.platforms.push(result[0])
  })
}

// 导入平台
function importPlatform() {
  OpenFile().then(response => {
    if (response.StatusCode == 200 && response.Data !== "") {
      ImportPlatform(response.Data).then(result => {
        loadPlatforms()
        MessagePlugin.success("导入成功")
      })
    }
    
  })
}

// 导出平台
function exportPlatform(platform) {
  SaveFile("", platform.Name+".json", "导出平台至").then((response)=>{
    console.log("savePath:",savePath)
    if (response.StatusCode == 200 && response.Data !== "") {
      var savePath = response.Data
      ExportPlatform(savePath, platform).then(result => {
        MessagePlugin.success("导出成功: "+ savePath)
      })
    }
  })
}

// 编辑平台
function editPlatform(platform) {
  GetPlatform(platform).then(result => {
    console.log("editPlatform->", result)
    interfacesStore.title = platform.Name
    interfacesStore.platform = result
    interfacesStore.visible = true
  })
}

// 更新平台
function updatePlatform(platform) {
  platform.Disabled = !platform.Disabled
  UpdatePlatform(platform).then(result => {
    loadPlatforms()  // 更新平台
    MessagePlugin.success('更新成功')
  })
}

// 删除平台
function deletePlatform(index, platform) {
  DeletePlatforms([platform]).then(result => {
    platformStore.platforms.splice(index, 1)
    MessagePlugin.success('删除成功')
  })
}

onMounted(() => {
  loadPlatforms()
})

onActivated(()=>{
  loadPlatforms()
})

</script>
<style scoped>
.box-card {
  margin-bottom: 10px;
}
</style>