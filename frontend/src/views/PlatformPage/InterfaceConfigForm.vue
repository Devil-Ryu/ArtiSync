<template>
  <t-form 
    label-align="top" 
    :rules="rules"
    :data="props.interfaceForm"
    scroll-to-first-error="smooth"
    @submit="onSubmit">
    <t-tabs default-value="request" placement="top">
      <t-tab-panel label="请求配置" value="request">
        <t-form-item name="prior" label="优先级" help="优先级决定接口执行的先后顺序">
          <t-input-number v-model="props.interfaceForm.Prior" theme="column" placeholder="请输入接口序号"></t-input-number>
        </t-form-item>
        <t-form-item name="serial" label="接口编码" help="统一平台内接口唯一标识">
          <t-input v-model="props.interfaceForm.Serial" disabled placeholder="请输入接口编码" />
        </t-form-item>
        <t-form-item name="alias" label="接口名称" help="接口别名，方便区分接口">
          <t-input v-model="props.interfaceForm.Name" placeholder="请输入接口名称" />
        </t-form-item>
        <t-form-item name="type" label="接口类型" help="接口类型，当为图片上传接口时，会循环遍历文章中的图片调用该接口或者组">
          <t-select v-model:value="props.interfaceForm.Type" :options="interfaceTypeOptions" />
        </t-form-item>
        <t-form-item name="url" label="接口URL" help="接口地址，不要带请求参数">
          <t-input v-model="props.interfaceForm.RequestURL" placeholder="请输入接口URL" />
        </t-form-item>
        <t-form-item label="请求方法">
          <t-radio-group v-model="props.interfaceForm.RequestMethod" size="small">
            <t-space :size="10">
              <t-radio-button value="GET">GET</t-radio-button>
              <t-radio-button value="POST">POST</t-radio-button>
              <t-radio-button value="PUT">PUT</t-radio-button>
              <t-radio-button value="OPTIONS">OPTIONS</t-radio-button>
              <t-radio-button value="PATCH">PATCH</t-radio-button>
              <t-radio-button value="DELETE">DELETE</t-radio-button>
              <t-radio-button value="HEAD">HEAD</t-radio-button>
              <t-radio-button value="TRACE">TRACE</t-radio-button>
              <t-radio-button value="CONNECT">CONNECT</t-radio-button>
            </t-space>
          </t-radio-group>
        </t-form-item>
        <t-divider />
        <t-form-item label="请求路径">
          <div style="width: 100%;">
            <t-row v-for="(item, index) in props.interfaceForm.RequestURLPath" :gutter="5" style="margin-top: 4px;">
                <t-col flex="auto">
                  <DynamicValue  v-model:item="props.interfaceForm.RequestURLPath[index]" v-model:cascaderOptions="interfacesStore.options"/>
                </t-col>
                <t-col flex="none">
                  <DeleteIcon size="15px" style="color: indianred"
                    @click="arrDelete(props.interfaceForm.RequestURLPath, index)" />
                </t-col>
            </t-row>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="arrAdd(props.interfaceForm.RequestURLPath)">添加查询参数</t-link>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="openBatchDialog(props.interfaceForm.RequestURLPath)">批量导入</t-link>
          </div>
        </t-form-item>
        <t-form-item label="查询参数">
          <div style="width: 100%;">
            <t-row v-for="(item, index) in props.interfaceForm.RequestQuery" :gutter="5" style="margin-top: 4px;">
                <t-col class="input-list-key">
                  <t-input 
                    v-model="item.Key"
                    placeholder="请输入键"
                    :status="replaceFilter(props.interfaceForm.Serial, 'Query', item.Key)?linkTheme:defaultTheme"
                  />
                </t-col>
                <t-col flex="auto">
                  <DynamicValue 
                    v-model:item="props.interfaceForm.RequestQuery[index]"
                    v-model:cascaderOptions="interfacesStore.options"
                    :status="replaceFilter(props.interfaceForm.Serial, 'Query', item.Key)?linkTheme:defaultTheme"
                    :disabled="replaceFilter(props.interfaceForm.Serial, 'Query', item.Key)"
                  />
                </t-col>
                <t-col flex="none">
                  <DeleteIcon size="15px" style="color: indianred"
                    @click="arrDelete(props.interfaceForm.RequestQuery, index)" />
                </t-col>
            </t-row>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="arrAdd(props.interfaceForm.RequestQuery)">添加查询参数</t-link>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="openBatchDialog(props.interfaceForm.RequestQuery)">批量导入</t-link>
          </div>

        </t-form-item>
        <t-form-item label="请求头">
          <div style="width: 100%;">
            <t-row v-for="(item, index) in props.interfaceForm.RequestHeaders" :gutter="5" style="margin-top: 4px;">
                <t-col class="input-list-key">
                  <t-input  
                    v-model="item.Key" 
                    placeholder="请输入键"  
                    :status="replaceFilter(props.interfaceForm.Serial, 'Headers', item.Key)?linkTheme:defaultTheme"
                  />
                </t-col>
                <t-col flex="auto">
                  <DynamicValue 
                    v-model:item="props.interfaceForm.RequestHeaders[index]"
                    v-model:cascaderOptions="interfacesStore.options"
                    :status="replaceFilter(props.interfaceForm.Serial, 'Headers', item.Key)?linkTheme:defaultTheme"
                    :disabled="replaceFilter(props.interfaceForm.Serial, 'Headers', item.Key)"
                    />
                    
                </t-col>
                <t-col flex="none">
                  <DeleteIcon size="15px" style="color: indianred"
                    @click="arrDelete(props.interfaceForm.RequestHeaders, index)" />
                </t-col>
            </t-row>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="arrAdd(props.interfaceForm.RequestHeaders)">添加请求头</t-link>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="openBatchDialog(props.interfaceForm.RequestHeaders)">批量导入</t-link>
          </div>
        </t-form-item>
        <t-form-item label="请求体" >
        <!-- <t-form-item label="请求体" v-if="['POST', 'PUT'].indexOf(props.interfaceForm.RequestMethod) !== -1"> -->
          <div style="width: 100%;">
            <t-row :gutter="5">
              <t-col  class="input-list-key">
                <t-select 
                  :options="['JSON', 'ROWDATA', 'FORMDATA', 'NONE'].map((value) => ({ label: value, value }))"
                  v-model="props.interfaceForm.RequestBodyType" />
              </t-col>
            </t-row>
            <t-row v-for="(item, index) in props.interfaceForm.RequestBody" :gutter="5" style="margin-top: 4px;">
                <t-col class="input-list-key">
                  <t-input
                    v-model="item.Key"
                    placeholder="请输入键"
                    :status="replaceFilter(props.interfaceForm.Serial, 'Body', item.Key)?linkTheme:defaultTheme"
                  />
                </t-col>
                <t-col flex="auto">
                  <DynamicValue
                    v-model:item="props.interfaceForm.RequestBody[index]"
                    v-model:cascaderOptions="interfacesStore.options"
                    :status="replaceFilter(props.interfaceForm.Serial, 'Body', item.Key)?linkTheme:defaultTheme"
                    :disabled="replaceFilter(props.interfaceForm.Serial, 'Body', item.Key)"
                  />
                </t-col>
                <t-col flex="auto"  v-if="item.IsFile">
                  <t-input v-model="item.FileName" placeholder="文件名称" />
                </t-col>
                <t-col flex="none">
                  <t-space :size="5">
                    <FileAttachmentIcon size="15px"  :style="{color: item.IsFile ? 'cornflowerblue':'gray', cursor: 'pointer'}" @click="item.IsFile = !item.IsFile"/>
                    <t-tooltip content="请求时转换字符串为对应格式数据">
                      <DataDisplayIcon size="15px" :style="{color: item.Convert ? 'cornflowerblue':'gray', cursor: 'pointer'}" @click="item.Convert = !item.Convert"/>
                    </t-tooltip>
                    <DeleteIcon size="15px" style="color: indianred;cursor: pointer"
                      @click="arrDelete(props.interfaceForm.RequestBody, index)" />
                  </t-space>
                </t-col>
            </t-row>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="arrAdd(props.interfaceForm.RequestBody)">添加请求体</t-link>
            <t-link size="small" style="margin-left: 4px;" theme="primary"
              @click="openBatchDialog(props.interfaceForm.RequestBody)">批量导入</t-link>

          </div>
        </t-form-item>
      </t-tab-panel>
      <t-tab-panel label="响应配置" value="response">
        <t-form-item label="响应体">
          <div style="width: 100%;">
            <t-row :gutter="5">
              <t-col class="input-list-key" >
                <t-select v-model:value="props.interfaceForm.ResponseType" :options="responseTypeOptions" />
              </t-col>
            </t-row>
            <t-textarea v-model:value="props.interfaceForm.ResponseTemple" :autosize="{ minRows: 5, maxRows: 10 }" />
          </div>
        </t-form-item>
        <t-form-item label="响应体取值">
          <div style="width: 100%;">
            <div style="margin-bottom: 10px; ">
              <t-radio-group  style="width: 100%;" v-model:value="props.interfaceForm.ResponseValidate"  size="small">
                <t-radio-button :value="true">进行响应体取值校验</t-radio-button>
                <t-radio-button :value="false">不进行响应体取值校验</t-radio-button>
              </t-radio-group>
            </div>
            <div class="input-list">
              <t-select class="input-list-key" v-model:value="props.interfaceForm.ResponseMapType"
                :options="responseMapTypeOptions" />
              <t-input class="input-list-value" style="width: 100%;" v-if="props.interfaceForm.ResponseMapType === 'NONE'" v-model="props.interfaceForm.ResponseMapRule"
                placeholder="请输入值" disabled />
              <t-input class="input-list-value" style="width: 100%;"  v-if="props.interfaceForm.ResponseMapType === 'RE'" v-model="props.interfaceForm.ResponseMapRule"
                placeholder="请输入值" />
              <t-cascader class="input-list-value" style="width: 100%;"  v-if="props.interfaceForm.ResponseMapType === 'JSON'" v-model="props.interfaceForm.ResponseMapRule"
                :options="interfacesStore.options" check-strictly />
            </div>
          </div>
          
        </t-form-item>
      </t-tab-panel>
    </t-tabs>
  </t-form>

  <BatchImportDialog />
</template>

<script setup>
import BatchImportDialog from "@/src/views/Components/BatchImportDialog.vue";
import { DeleteIcon, DataDisplayIcon, FileAttachmentIcon} from "tdesign-icons-vue-next";
import { computed, ref } from "vue";
import { useInterfacesStore, useBatchImportStore } from "@/src/store/platform"
import { MessagePlugin } from "tdesign-vue-next";
import DynamicValue from "@/src/views/Components/DynamicValue.vue";
const interfacesStore = useInterfacesStore()
const batchImportStore = useBatchImportStore()

const linkTheme = "success"
const defaultTheme = "default"

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
        Serial: undefined,
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
        ResponseValidate: false,
      }
    }
  },
})
const emits = defineEmits(['close', 'save'])

const rules = {
    prior: [ {required : true}],
    serial: [ {required : true}],
    alias: [{required: true},
      {max: 5, message: '<20>', trigger: 'change'}
    ],
    type: [ {required : true}],
    url: [ {required : true}],
}

const interfaceTypeOptions = [
  { label: "普通接口", value: "normal" },
  { label: "图片上传接口", value: "images" },
  { label: "文章上传接口", value: "article" },
]
const responseTypeOptions = [
  { label: 'NONE', value: 'NONE' },
  { label: 'JSON', value: 'JSON' },
  { label: 'FORM', value: 'FORM' },
  { label: 'ROWDATA', value: 'ROWDATA' },
]

const responseMapTypeOptions = [
  { label: 'JSON', value: 'JSON' },
  { label: 'NONE', value: 'NONE' },
  { label: 'RE', value: 'RE' },
]

// const batchImportDialogVisible = ref(false)
// const batchImportSelectOptions = [
//   { label: 'JSON格式', value: 'JSON' },
//   { label: 'ROWDATA格式', value: 'ROWDATA' },
//   // { label: 'SQLITE文件格式', value: 'SQLITE' },
// ]
// const batchImportContentPlaceholder =
// {
//   "JSON": "{'key':'value'}",
//   "ROWDATA": "key1=value1\nkey2=value2",
//   "SQLITE": "Sqlite文件",
// }
// const batchImportDialogSelectValue = ref("JSON")
// const batchImportDialogContent = ref("")
// const batchImportArr = ref([])


function arrAdd(arr) {
  arr.push({ Dynamic: false })
}

function arrDelete(arr, index) {
  arr.splice(index, 1)[0];
}

function openBatchDialog(arr) {
  batchImportStore.typeInsert = false
  batchImportStore.contentArr = arr
  batchImportStore.visible = true
}

const onSubmit = ({ validateResult, firstError, e }) => {
  e.preventDefault();
  if (validateResult === true) {
    MessagePlugin.success('提交成功');
  } else {
    console.log('Validate Errors: ', firstError, validateResult);
    MessagePlugin.warning(firstError);
  }
};

function replaceFilter(Serail, Type, Key) {
  var result = interfacesStore.platform.ReplaceMaps.filter((item)=>{return item.Interfaces.indexOf(Serail)!==-1 && item.Type===Type && item.Key===Key && !item.Disabled})
  return result.length !== 0
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
}

.input-list-value {
  width: 80%;
  margin-right: 10px;
}
</style>