<template>
  <t-dialog :visible="interfacesStore.visible" :header="interfacesStore.title" destroyOnClose cancel-btn="取消"
    confirm-btn="保存" size="100%" width="90%" @cancel="clickBtn = 'cancel';" @close="onClose" @confirm="saveInterfaces"
    @opened="onOpen">
    <template #header>
      <t-row justify="space-between" style="width: 98%;">
        <t-col>
          <t-input v-model:value="interfacesStore.platform.Name" style="width: 50%;" />
        </t-col>
      </t-row>
    </template>
    <template #body>
      <t-tabs default-value="interface" placement="top">
        <t-tab-panel label="编辑接口" value="interface">
          <t-collapse expand-icon-placement="right" expand-mutex borderless>
            <div v-for="(interface_, index) in interfacesStore.platform.Interfaces">
              <t-card v-if="interface_.ParentSerial == undefined" :shadow="true" style="margin-top: 10px; padding: 0">
                <t-collapse-panel destroyOnCollapse>
                  <template #header>
                    <t-tag size="medium" shape="mark" theme="primary" variant="outline">优先级: {{ interface_.Prior
                    }}</t-tag>
                    <div style="margin-left: 10px;"> {{ interface_.Name }} </div>
                  </template>
                  <template #headerRightContent>
                    <ChevronUpIcon @click="sortInterface('up', interface_)" size="18" style="cursor: pointer" />
                    <ChevronDownIcon @click="sortInterface('down', interface_)" size="18" style="cursor: pointer" />
                    <DeleteIcon @click="deleteInterfaces(index, interface_)" style="cursor: pointer" />
                  </template>
                  <InterfaceConfigForm v-if="!interfacesStore.platform.Interfaces[index].IsGroup"
                    v-model:interface-form="interfacesStore.platform.Interfaces[index]" />

                  <InterfaceGroupForm v-if="interfacesStore.platform.Interfaces[index].IsGroup"
                    v-model:interface-form="interfacesStore.platform.Interfaces[index]" />
                </t-collapse-panel>
              </t-card>
            </div>
          </t-collapse>
          <t-dropdown :options="[{ content: '新增接口', value: false }, { content: '新增接口组', value: true }]" placement="right"
            @click="addInterface">
            <t-button shape="circle" theme="primary" style="margin-top: 10px;">
              <template #icon>
                <AddIcon />
              </template>
            </t-button>
          </t-dropdown>
        </t-tab-panel>
        <t-tab-panel label="编辑规则" value="mapRules">
          <div style="margin-top: 8px;">
            <t-row v-for="(item, index) in interfacesStore.platform.ReplaceMaps" :gutter="5" style="margin-top: 4px;" align="center" >
              <t-col>
                <CheckSelect v-model:item="interfacesStore.platform.ReplaceMaps[index]" />
              </t-col>
              <t-col style="width: 15%;" >
                <t-input placeholder="请输入键" v-model="item.Key" />
              </t-col>
              <t-col style="width: 25%;" >
                <t-input placeholder="请输入值" v-model="item.Value" />
              </t-col>
              <t-col style="width: 40%;" >
                <t-select placeholder="请选择接口" v-model:value="item.Interfaces" :options="interfacesStore.SerialOptions"
                  multiple :min-collapsed-num="2" />
              </t-col>
              <t-col style="width: 4%; text-align: center;" >
                <DeleteIcon size="15px" style="color: indianred" 
                @click="arrDelete(interfacesStore.platform.ReplaceMaps, index)"
                />
              </t-col>
            </t-row>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="arrAdd(interfacesStore.platform.ReplaceMaps)">添加规则</t-link>
            <!-- <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="openBatchDialog(interfacesStore.platform.ReplaceMaps)">批量导入</t-link> -->
          </div>
        </t-tab-panel>
      </t-tabs>
    </template>
  </t-dialog>
</template>

<script lang="jsx" setup>
import { DialogPlugin, Input, MessagePlugin } from "tdesign-vue-next";
import InterfaceGroupForm from "./InterfaceGroupForm.vue";
import InterfaceConfigForm from "./InterfaceConfigForm.vue";
import { CreateInterface, DeleteInterfaces, UpdateInterface, UpdatePlatform } from "../../../wailsjs/go/api/DBController.js";
import { AddCircleIcon, AddIcon, DeleteIcon, ChevronUpIcon, ChevronDownIcon, CheckRectangleFilledIcon } from "tdesign-icons-vue-next";
import { ref, watch } from "vue";
import { useInterfacesStore, usePlatformStore } from "@/src/store/platform"
import CheckSelect from "./Components/CheckSelect.vue";

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
      theme: "warning",
      confirmBtn: "确认",
      cancelBtn: "取消",
      closeBtn: false,
      onConfirm: ({ e }) => { confirmDialog.hide(); interfacesStore.visible = false; clickBtn.value = "" },
      onCancel: () => { confirmDialog.hide(); clickBtn.value = "" },
    })
  } else {
    interfacesStore.visible = false
    clickBtn.value = ""
  }
}

// 接口排序
function sortInterface(method, interface_) {
  if (method === "up") {  // 向上则序号减一
    interface_.Prior = parseInt(interface_.Prior) - 1
  } else { // 向上则序号加一
    interface_.Prior = parseInt(interface_.Prior) + 1
  }
  interfacesStore.platform.Interfaces.sort(function (a, b) { return parseInt(a.Prior) >= parseInt(b.Prior) })
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
    var platformIndex = platformStore.platforms.findIndex((item) => { return item.ID == platform.ID })
    platformStore.platforms[platformIndex] = platform
  })
}

function arrAdd(arr) {
  arr.push({ disabled: false })
}

function arrDelete(arr, index) {
  arr.splice(index, 1)[0];
}

// function openBatchDialog(arr) {
//   batchImportDialogVisible.value = true
//   batchImportArr.value = arr
// }

</script>

<style scoped></style>