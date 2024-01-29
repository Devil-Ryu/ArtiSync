<template>
  <t-dialog
      :visible="interfacesStore.visible"
      :header="interfacesStore.title"
      destroyOnClose
      cancel-btn="取消"
      confirm-btn="保存"
      size="100%"
      width="90%"
      @cancel="clickBtn = 'cancel';"
      @close="onClose"
      @confirm="saveInterfaces"
      @opened="onOpen"
  >
    <template #header>
      <t-input  v-model:value="interfacesStore.platform.Name" style="width: 50%;" />
    </template>
    <template #body>
      <div v-for="(interface_, index) in interfacesStore.platform.Interfaces">
        <t-card  v-if="interface_.ParentSerial == undefined" :shadow="true" style="margin-top: 10px; padding: 0">
          <t-collapse expand-icon-placement="right" expand-mutex borderless >
            <t-collapse-panel destroyOnCollapse>
              <template #header>
                <t-tag size="medium" shape="mark" theme="primary" variant="outline">优先级: {{interface_.Prior}}</t-tag>
                <div style="margin-left: 10px;"> {{ interface_.Name }} </div>
              </template>
              <template #headerRightContent>
                <ChevronUpIcon @click="sortInterface('up', interface_)" size="18" style="cursor: pointer"/>
                <ChevronDownIcon @click="sortInterface('down', interface_)" size="18" style="cursor: pointer"/>
                <DeleteIcon @click="deleteInterfaces(index,interface_)" style="cursor: pointer"/>
              </template>
              <InterfaceConfigForm
                  v-if="!interfacesStore.platform.Interfaces[index].IsGroup"
                  v-model:interface-form="interfacesStore.platform.Interfaces[index]"/>

              <InterfaceGroupForm
                  v-if="interfacesStore.platform.Interfaces[index].IsGroup"
                  v-model:interface-form="interfacesStore.platform.Interfaces[index]"/>
            </t-collapse-panel>
          </t-collapse>
        </t-card>
      </div>
      <t-dropdown
        :options="[{ content: '新增接口', value: false }, { content: '新增接口组', value: true }]"
        placement="right"
        @click="addInterface"
      >
      <t-button shape="circle" theme="primary"  style="margin-top: 10px;" >
          <template #icon><AddIcon /></template>
        </t-button>
    </t-dropdown>
    </template>
  </t-dialog>
</template>

<script lang="jsx" setup>
import {DialogPlugin, Input, MessagePlugin} from "tdesign-vue-next";
import InterfaceGroupForm from "./InterfaceGroupForm.vue";
import InterfaceConfigForm from "./InterfaceConfigForm.vue";
import {CreateInterface, DeleteInterfaces, UpdateInterface, UpdatePlatform} from "../../../wailsjs/go/api/DBController.js";
import {AddCircleIcon, AddIcon, DeleteIcon, ChevronUpIcon, ChevronDownIcon} from "tdesign-icons-vue-next";
import {ref, watch} from "vue";
import { useInterfacesStore, usePlatformStore } from "@/src/store/platform"
const interfacesStore = useInterfacesStore()
const platformStore = usePlatformStore()

const clickBtn = ref("")
let oldForm = undefined

function onOpen() {
  oldForm = JSON.stringify(interfacesStore.platform);
}

function onClose() {
  if (clickBtn.value === "cancel") {
    interfacesStore.visible = false
    clickBtn.value = ""
    return
  }
  if (JSON.stringify(interfacesStore.platform) !== oldForm) {
    const confirmDialog = DialogPlugin.confirm({
      header: "警告",
      body: "改动尚未保存，确认退出么?",
      theme:"warning",
      confirmBtn: "确认",
      cancelBtn: "取消",
      closeBtn: false,
      onConfirm: ({e}) => {confirmDialog.hide(); interfacesStore.visible = false;clickBtn.value = ""},
      onCancel: () => {confirmDialog.hide();clickBtn.value = ""},
    })
  } else {
    interfacesStore.visible = false
    clickBtn.value = ""
  }
}

// 接口排序
function sortInterface(method, interface_) {
  if (method === "up") {  // 向上则序号减一
    interface_.Prior =  parseInt(interface_.Prior) - 1
  } else { // 向上则序号加一
    interface_.Prior =  parseInt(interface_.Prior) + 1
  }
  interfacesStore.platform.Interfaces.sort(function(a, b){return parseInt(a.Prior) >= parseInt(b.Prior)})
}

// 添加新接口
function addInterface(data) {
  var tmp = {
    platformID: interfacesStore.platform.ID,
    Prior: 1,
    IsGroup: data.value,
  }
  CreateInterface(tmp).then(result => {
    interfacesStore.platform.Interfaces.push(result)
  })
}

// 保存接口信息
function saveInterfaces() {
  console.log("save interface", interfacesStore.platform.Interfaces)
  UpdateInterface(interfacesStore.platform.Interfaces).then(result => {
    MessagePlugin.success('保存成功')
  })
  updatePlatform(interfacesStore.platform)
  interfacesStore.visible = false
}

// 删除接口信息
function deleteInterfaces(index, interface_) {
  console.log("in deleteInterfaces")
  DeleteInterfaces([interface_]).then(result => {
    interfacesStore.platform.Interfaces.splice(index, 1)
    MessagePlugin.success('删除成功')
  })

}

// 更新平台
function updatePlatform(platform) {
  UpdatePlatform(platform).then(result => {
    var platformIndex = platformStore.platforms.findIndex((item)=> {return item.ID == platform.ID})
    platformStore.platforms[platformIndex] = platform
  })
}
</script>

<style scoped>


</style>