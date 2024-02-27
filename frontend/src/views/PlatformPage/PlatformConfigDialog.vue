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
      <t-tabs default-value="mapRules" placement="top" @change="onTabChange">
        <t-tab-panel label="编辑规则" value="mapRules">
          <div style="margin-top: 8px;">
            <t-row v-for="(item, index) in interfacesStore.platform.ReplaceMaps" :gutter="5" style="margin-top: 4px;"
              align="center">
              <t-col>
                <CheckSelect v-model:item="interfacesStore.platform.ReplaceMaps[index]" />
              </t-col>
              <t-col style="width: 15%;">
                <t-input placeholder="请输入键" v-model="item.Key" />
              </t-col>
              <t-col style="width: 25%;">
                <t-input placeholder="请输入值" v-model="item.Value" />
              </t-col>
              <t-col style="width: 38%;">
                <t-select placeholder="请选择接口" v-model:value="item.Interfaces" :options="interfacesStore.SerialOptions"
                  multiple :min-collapsed-num="1">
                  <template #panelTopContent>
                    <t-row justify="space-around">
                      <t-col :span="6">
                        <t-button theme="primary" variant="text" block @click="selectAll(item)">
                          全选
                        </t-button>
                      </t-col>
                      <t-col :span="6">
                        <t-button theme="primary" variant="text" block @click="unSelectAll(item)">
                          取消全选
                        </t-button>
                      </t-col>
                    </t-row>
                  </template>
                </t-select>
              </t-col>
              <t-col style="width: 5%; text-align: center;">
                <DeleteIcon size="15px" style="color: indianred"
                  @click="arrDelete(interfacesStore.platform.ReplaceMaps, index)" />
              </t-col>
            </t-row>
            <t-row justify="space-between">
              <t-col>
                <t-link size="small" style="margin-left: 4px;" theme="primary"
                  @click="arrAdd(interfacesStore.platform.ReplaceMaps)">添加规则</t-link>
                <t-link size="small" style="margin-left: 4px;" theme="primary"
                  @click="openBatchDialog(interfacesStore.platform.ReplaceMaps)">批量导入</t-link>
                <t-link size="small" style="margin-left: 4px;" theme="primary"
                  @click="ruleFullStatus(false)">全部启用</t-link>
                <t-link size="small" style="margin-left: 4px;" theme="danger" @click="ruleFullStatus(true)">全部禁用</t-link>
              </t-col>
            </t-row>
          </div>
        </t-tab-panel>
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
                    <t-tooltip content="运行">
                      <PlayIcon @click="runInterface(interface_)" size="18"
                        style="cursor: pointer; color: var(--td-brand-color-7);" />
                    </t-tooltip>
                    <t-tooltip content="上移">
                      <ChevronUpIcon @click="sortInterface('up', interface_)" size="18" style="cursor: pointer" />
                    </t-tooltip>
                    <t-tooltip content="下移">
                      <ChevronDownIcon @click="sortInterface('down', interface_)" size="18" style="cursor: pointer" />
                    </t-tooltip>
                    <t-tooltip content="删除">
                      <DeleteIcon @click="deleteInterfaces(index, interface_)"
                        style="cursor: pointer; color: var(--td-error-color);" />
                    </t-tooltip>
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
        <t-tab-panel label="运行记录" value="records"> <InterfaceRecordPage /></t-tab-panel>
        <t-tab-panel label="测试缓存" value="caches"><CacheView /> </t-tab-panel>
      </t-tabs>
    </template>
  </t-dialog>
  <BatchImportDialog />
</template>

<script lang="jsx" setup>
import { DialogPlugin, MessagePlugin } from "tdesign-vue-next";
import InterfaceGroupForm from "./InterfaceGroupForm.vue";
import InterfaceConfigForm from "./InterfaceConfigForm.vue";
import { CreateInterface, DeleteInterfaces, UpdateInterface, UpdatePlatform } from "../../../wailsjs/go/api/DBController.js";
import {  AddIcon, DeleteIcon, ChevronUpIcon, ChevronDownIcon, PlayIcon } from "tdesign-icons-vue-next";
import { ref, watch } from "vue";
import { useInterfacesStore, usePlatformStore, useInterfaceRecordsStore, useBatchImportStore } from "@/src/store/platform"
import BatchImportDialog from "../Components/BatchImportDialog.vue";
import CheckSelect from "./components/CheckSelect.vue";
import CacheView from "./components/CacheView.vue"
import { EventsOff, EventsOn } from "@/wailsjs/runtime/runtime";
import { TestInterface, DeleteTestNetControllerCache, GetTestNetControllerInfo} from "@/wailsjs/go/api/ATController";
import InterfaceRecordPage from "../InterfaceRecords/InterfaceRecordPage.vue";

const batchImportStore = useBatchImportStore()
const interfacesStore = useInterfacesStore()
const platformStore = usePlatformStore()
const interfaceRecordStore = useInterfaceRecordsStore()

const clickBtn = ref("")
let oldForm = undefined

watch(
  () => interfacesStore.platform.ReplaceMaps,
  (oldValue, newValue) => {
    replaceMapConvert()
  },
  { deep: true }
)

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

function onTabChange(tabValue) {
  console.log("当前标签[tab]: ", tabValue)
  if (tabValue === "caches") {
    GetTestNetControllerInfo().then((result)=>{
      interfaceRecordStore.testCaches = result
    }).catch((err)=>{
      MessagePlugin.error("获取缓存失败: "+err)
    })
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
  // 更新接口信息
  UpdateInterface(interfacesStore.platform.Interfaces).then(result => {
    MessagePlugin.success('保存成功')
  })
  updatePlatform(interfacesStore.platform)  // 更新平台信息
  oldForm = JSON.stringify(interfacesStore.platform);  // 更新旧表单
  // interfacesStore.visible = false
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
  arr.push({})
}

function arrDelete(arr, index) {
  arr.splice(index, 1)[0];
}

function openBatchDialog(arr) {
  batchImportStore.typeInsert = true
  batchImportStore.contentArr = arr
  batchImportStore.visible = true
}

function selectAll(item) {
  item.Interfaces = interfacesStore.SerialOptions.map(item => { return item.value })
}

function unSelectAll(item) {
  item.Interfaces = []
}

function ruleFullStatus(disabled) {
  for (let i = 0; i < interfacesStore.platform.ReplaceMaps.length; i++) {
    interfacesStore.platform.ReplaceMaps[i].Disabled = disabled
  }
  if (disabled === true) {
    MessagePlugin.error("全部禁用成功")
  } else {

    MessagePlugin.success("全部启用成功")
  }
}

function replaceMapConvert() {
  interfacesStore.platform.ReplaceMaps.forEach((ruleItem) => {
    // 判断规则是否被禁用
    if (ruleItem.Interfaces && !ruleItem.Disabled) {
      // 遍历规则应用的接口编号
      ruleItem.Interfaces.forEach((interfaceSerial) => {
        replaceInterfaceMap(interfaceSerial, ruleItem.Type, ruleItem.Key, ruleItem.Value, interfacesStore.platform.Interfaces)
      })
    }
  })
}

// 找到接口并替换对应区域
function replaceInterfaceMap(Serial, Type, Key, Value, interfaceList) {
  // 遍历接口列表
  for (let i = 0; i < interfaceList.length; i++) {
    // 找到要替换的接口
    if (interfaceList[i].Serial === Serial) {
      // 找到替换的区域
      if (Type === "Query") {
        replaceInterfaceAttr(interfaceList[i].RequestQuery, Key, Value)
      }
      if (Type === "Headers") {
        replaceInterfaceAttr(interfaceList[i].RequestHeaders, Key, Value)
      }
      if (Type === "Body") {
        replaceInterfaceAttr(interfaceList[i].RequestBody, Key, Value)
      }
    }

    // 如果为组则在组内搜索接口
    if (interfaceList[i].IsGroup === true) {
      replaceInterfaceMap(Serial, Type, Key, Value, interfaceList[i].Children)
    }

  }
}

// 替换接口对应区域
function replaceInterfaceAttr(arr, Key, Value) {
  var keyIndex = arr.findIndex((item) => item.Key === Key)
  if (keyIndex !== -1) {
    arr[keyIndex]["Value"] = Value
  } else {
    arr.push({ Key: Key, Value: Value })
  }
}

function runInterface(interfaceInfo) {
  MessagePlugin.info("运行接口: "+interfaceInfo.Name)
  TestInterface(interfacesStore.platform, interfaceInfo).then(() => {
    MessagePlugin.success("接口运行完毕: "+interfaceInfo.Name)
  }).catch((err) => {
    MessagePlugin.error("运行失败："+err)
  })
}



</script>

<style scoped></style>