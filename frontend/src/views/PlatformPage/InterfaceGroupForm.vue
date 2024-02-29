<template>
  <t-form label-align="top" show-error-message>
    <t-form-item label="优先级">
      <t-input-number v-model="props.interfaceForm.Prior" theme="column" placeholder="请输入接口序号"></t-input-number>
    </t-form-item>
    <t-form-item label="接口组名称">
      <t-input v-model="props.interfaceForm.Name" placeholder="请输入接口名称" />
    </t-form-item>
    <t-form-item label="接口组类型">
      <t-select v-model:value="props.interfaceForm.Type" :options="interfaceTypeOptions" />
    </t-form-item>
    <t-divider />
    <t-form-item label="接口列表">
      <div style="width: 100%;">
        <t-collapse expand-icon-placement="right" expand-mutex borderless>
          <div class="input-list" v-for="(interface_, index) in props.interfaceForm.Children">
            <t-card :shadow="true" style="width: 100%;">
              <t-collapse-panel destroyOnCollapse :disabled="interface_.Disabled">
                <template #header>
                  <t-tag size="medium" shape="mark" theme="primary" variant="outline">优先级: {{ interface_.Prior }}</t-tag>
                  <div style="margin-left: 10px;"> {{ interface_.Name }} </div>
                </template>
                <template #headerRightContent>
                  <t-tooltip content="点击启用" v-if="interface_.Disabled">
                      <LockOnIcon @click="changeInterfaceStatus(interface_)" size="16" style="cursor: pointer; "/>
                    </t-tooltip>
                    <t-tooltip content="点击禁用" v-if="!interface_.Disabled">
                      <LockOffIcon @click="changeInterfaceStatus(interface_)" size="16" style="cursor: pointer;"/>
                    </t-tooltip>
                    <t-tooltip content="运行">
                      <PlayIcon @click="runInterface(interface_)" size="22"
                        style="cursor: pointer; color: var(--td-brand-color-7);" />
                    </t-tooltip>
                    <t-tooltip content="上移">
                      <ChevronUpIcon @click="sortInterface('up', interface_)" size="22" style="cursor: pointer" />
                    </t-tooltip>
                    <t-tooltip content="下移">
                      <ChevronDownIcon @click="sortInterface('down', interface_)" size="22" style="cursor: pointer" />
                    </t-tooltip>
                    <t-tooltip content="删除">
                      <DeleteIcon @click="deleteInterfaces(index, interface_)"
                        style="cursor: pointer; color: var(--td-error-color);" />
                    </t-tooltip>
                </template>
                <InterfaceConfigForm v-model:interface-form="props.interfaceForm.Children[index]" />
              </t-collapse-panel>
            </t-card>
          </div>
        </t-collapse>

        <t-button theme="primary" variant="dashed" style="width: 100%; margin-top: 5px;"
          @click="addInterface">添加接口</t-button>
      </div>
    </t-form-item>
  </t-form>
</template>

<script setup>
import { MessagePlugin } from "tdesign-vue-next";
import InterfaceConfigForm from "./InterfaceConfigForm.vue";
import { DeleteIcon, ChevronUpIcon, ChevronDownIcon, PlayIcon, LockOffIcon, LockOnIcon } from "tdesign-icons-vue-next";
import { CreateInterface, DeleteInterfaces, UpdateInterface } from "@/wailsjs/go/api/DBController.js";
import { TestInterface } from "@/wailsjs/go/api/ATController";
import { useInterfacesStore } from "@/src/store/platform"
const interfacesStore = useInterfacesStore()

const props = defineProps({
  title: {
    type: String,
    default: '接口配置'
  },
  interfaceForm: {
    type: Object,
    default: () => {
      return {
        ID: undefined,
        PlatformID: undefined,
        Index: undefined,
        Name: "",
        IsGroup: true,
        ParentSerial: undefined,
        Children: [],
        RequestURL: "",
        RequestURLPath: [],
        RequestMethod: "",
        RequestQuery: [],
        RequestHeaders: [],
        RequestBody: [],
        ResponseType: "",
        ResponseTemple: "",
      }
    }
  },
})
const emits = defineEmits(['close', 'save'])

const interfaceTypeOptions = [
  { label: "普通接口", value: "normal" },
  { label: "图片上传接口", value: "images" },
  { label: "文章上传接口", value: "article" },
]

// 接口排序
function sortInterface(method, interface_) {
  if (method === "up") {  // 向上则序号减一
    interface_.Prior = parseInt(interface_.Prior) - 1
  } else { // 向上则序号加一
    interface_.Prior = parseInt(interface_.Prior) + 1
  }

  function compare(key) {
    return (a, b) => {
      let val1 = parseInt(a[key]);
      let val2 = parseInt(b[key]);
      return val1 >= val2
    }
  }

  props.interfaceForm.Children.sort(compare('Prior'))
}

// 添加新接口
function addInterface() {
  var tmp = {
    Prior: 1,
    IsGroup: false,
    ParentSerial: props.interfaceForm.Serial
  }
  CreateInterface(tmp).then(result => {
    props.interfaceForm.Children.push(result)
  })
}

// 删除接口信息
function deleteInterfaces(index, interface_) {
  console.log("in deleteInterfaces")
  DeleteInterfaces([interface_]).then(result => {
    props.interfaceForm.Children.splice(index, 1)
    MessagePlugin.success('删除成功')
  })

}

function runInterface(interfaceInfo) {
  MessagePlugin.info("运行接口: "+interfaceInfo.Name)
  TestInterface(interfacesStore.platform, interfaceInfo).then(() => {
    MessagePlugin.success("接口运行完毕: "+interfaceInfo.Name)
  }).catch((err) => {
    MessagePlugin.error("运行失败："+err)
  })
}

// changeInterfaceStatus 改变接口状态
function changeInterfaceStatus(interfaceInfo) {
  interfaceInfo.Disabled = !interfaceInfo.Disabled
  UpdateInterface([interfaceInfo]).then(result => {
    if (interfaceInfo.Disabled) {
      MessagePlugin.error('接口已禁用')
    } else {

      MessagePlugin.success('接口已启用')
    }
  })

}

</script>

<style scoped>
.input-list {
  display: flex;
  align-items: center;
  width: 100%;
  margin-bottom: 5px;
}

.input-list-key {
  width: 20%;
  margin-right: 10px;
}

.input-list-value {
  width: 60%;
  margin-right: 10px;
}

.input-swtich-type {
  width: 60px;
  margin-right: 10px;
}

.input-list-type {
  width: 15%;
  margin-right: 10px;
  margin-bottom: 10px;
}
</style>