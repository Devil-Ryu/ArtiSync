<template>
  <t-form label-align="top"  show-error-message>
    <t-form-item label="优先级">
      <t-input-number v-model="props.interfaceForm.Prior" theme="column"  placeholder="请输入接口序号"></t-input-number>
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
        <div class="input-list" v-for="(interface_, index) in props.interfaceForm.Children">
          <t-card  :shadow="true" style="width: 100%;">
            <t-collapse expand-icon-placement="right" borderless >
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
                    v-model:interface-form="props.interfaceForm.Children[index]"
                />
              </t-collapse-panel>
            </t-collapse>
          </t-card>
        </div>
        <t-button 
          theme="primary" 
          variant="dashed" 
          style="width: 100%; margin-top: 5px;"
          @click="addInterface">添加接口</t-button>
      </div>
    </t-form-item>
  </t-form>
</template>

<script setup>
import {MessagePlugin} from "tdesign-vue-next";
import InterfaceConfigForm from "./InterfaceConfigForm.vue";
import { DeleteIcon, ChevronUpIcon, ChevronDownIcon} from "tdesign-icons-vue-next";
import {CreateInterface, DeleteInterfaces} from "@/wailsjs/go/api/DBController.js";
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
  { label: "普通接口", value: "normal"},
  { label: "图片上传接口", value: "images"},
  { label: "文章上传接口", value: "article"},
]

// 接口排序
function sortInterface(method, interface_) {
  if (method === "up") {  // 向上则序号减一
    interface_.Prior =  parseInt(interface_.Prior) - 1
  } else { // 向上则序号加一
    interface_.Prior =  parseInt(interface_.Prior) + 1
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