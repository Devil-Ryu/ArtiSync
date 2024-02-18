<template>
  <t-drawer size="100%" :header="false" :footer="true" placement="left" :confirm-btn="null" cancel-btn="返回"
    v-model:visible="articleStore.articleDetailVisible">
    <div style="height: 400px; width: 100%; ">
      <t-card>
        <template #title>
          <div style="margin-left: 14px; ">基本信息</div>
        </template>
        <t-descriptions colon style="margin-top: -20px;" size="small" :label-style="{ width: '1%' }"
          :content-style="{ width: '5%' }" :column="2">

          <t-descriptions-item label="文章名称">{{ articleStore.articleDetailBasicInfo.Title }}</t-descriptions-item>
          <t-descriptions-item label="文章类型"><t-tag theme="warning">{{ articleStore.articleDetailBasicInfo.Type
          }}</t-tag></t-descriptions-item>

          <t-descriptions-item label="图片数量">{{ articleStore.articleDetailBasicInfo.ImageNum }}</t-descriptions-item>
          <t-descriptions-item label="接口数量">{{ articleStore.articleDetailBasicInfo.InterfaceNum }}</t-descriptions-item>

          <t-descriptions-item label="发布状态"><t-tag :theme="statusNameListMap[articleStore.articleDetailBasicInfo.Status].theme">{{ statusNameListMap[articleStore.articleDetailBasicInfo.Status].label
          }}</t-tag></t-descriptions-item>
          <t-descriptions-item label="发布平台">{{ articleStore.articleDetailBasicInfo.Platforms.join(', ')
          }}</t-descriptions-item>
          <t-descriptions-item label="记录ID">{{ articleStore.articleDetailBasicInfo.InterfaceRecordID }}</t-descriptions-item>

        </t-descriptions>
      </t-card>

      <!-- 平台上传进度 -->
      <t-card style="margin-top: 10px;">
        <template #title>
          <div style="margin-left: 14px; ">平台上传进度</div>
        </template>
        <t-row :gutter="[0, 24]" >
          <t-col :span="colSpan" v-for="platform in articleStore.articleDetailPlatformsList ">
            <t-progress theme="circle" :percentage="platform.RunProgress" >
            <!-- <t-progress theme="circle" :percentage="platform.RunProgress" :status="statusNameListMap[platform.RunStatus].theme"> -->
              <template #label>
                <div class="platform-box">
                  <div class="platform-name">{{ platform.Name }}</div>
                  <t-tag class="platform-tag" :theme="statusNameListMap[platform.RunStatus].theme">{{ statusNameListMap[platform.RunStatus].label }}</t-tag>
                  <div class="platform-progress"> {{ platform.InterfaceActualRunNum }} / {{ platform.InterfacePredictRunNum }} </div>
                </div>
              </template>
            </t-progress>
          </t-col>
        </t-row>

      </t-card>

      <!-- 接口运行状态 -->
      <t-card style="margin-top: 10px;">
        <template #title>
          <div style="margin-left: 14px; ">接口运行状态</div>
        </template>
        <t-form>
          <t-row>
            <t-col :span="10">
              <t-row :gutter="[10, 24]">
                <t-col :span="4">
                  <t-form-item label="平台名称" name="platforName">
                    <t-select multiple :min-collapsed-num="1" :options="[]" v-model="filterValue.level" />
                  </t-form-item>
                </t-col>
                <t-col :span="4">
                  <t-form-item label="接口名称" name="interface">
                    <t-input type="search" placeholder="请输入接口名称" v-model="filterValue.tag" />
                  </t-form-item>
                </t-col>
                <t-col :span="4">
                  <t-form-item label="运行状态" name="status">
                    <t-select multiple :min-collapsed-num="1" :options="[]" v-model="filterValue.level" />
                  </t-form-item>
                </t-col>
              </t-row>
            </t-col>
            <t-col :span="2" class="operation-container">
              <t-button theme="primary" type="submit" :style="{ marginLeft: '8px' }" @click="qeury"> 查询 </t-button>
              <t-button type="reset" variant="base" theme="default" @click="reset"> 重置 </t-button>
            </t-col>
          </t-row>
        </t-form>
        <t-table :columns="columns" :data="data" row-key="PlatformName">
          <template #Status="{ row }">
            <t-tag :theme="statusNameListMap[row.Status].theme">
              {{ statusNameListMap[row.Status].label }}
            </t-tag>
          </template>
          <template #operation="slotProps">
          <t-link theme="primary">详情</t-link>
        </template>
        </t-table>
      </t-card>

    </div>
  </t-drawer>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useArticleStore, statusNameListMap } from "@/src/store/article"
import { useInterfaceRecordsStore } from "@/src/store/platform";
import { EventsOn } from "@/wailsjs/runtime/runtime";
import { GetInterfaceRecords } from "@/wailsjs/go/api/DBController";
const articleStore = useArticleStore()
const interfaceRecordsStore = useInterfaceRecordsStore()

const colSpan = ref(2)
const filterValue = ref({});

const columns = ref([
  { colKey: 'PlatformName', title: '平台名称', width: 150, },
  { colKey: 'Name', title: '接口名称', width: 150, ellipsis: true, },
  { colKey: 'RequestURL', title: '接口URL', width: 240, ellipsis: true, },
  { colKey: 'ResponseMessage', title: '响应内容', width: 240, ellipsis: true, },
  { colKey: 'Status', title: '运行状态', width: 120, ellipsis: true, },
  { colKey: 'operation', title: '操作', width: 80, fixed: 'right', align: 'center'},
])


const data = ref([...interfaceRecordsStore.records])

watch(
  ()=>interfaceRecordsStore.records,
  ()=>{
    data.value = [...interfaceRecordsStore.records]
  }
)

// TODO(待添加获取接口记录的函数，以及删除旧的更新图片状态的逻辑)
// onMounted(()=>{
//   test()
// })

// 函数区
EventsOn("UpdateInterfaceRecord", (recordID)=>{
  GetInterfaceRecords({"record_id": recordID}).then(result=>{
    interfaceRecordsStore.records = result
    console.log("result: ", result)
  })
})

</script>

<style scoped>
.platform-box {
  .platform-name {
    margin-top: 14px;
    margin-left: 16px;
    margin-right: 16px;
    font-size: 18px;
    text-align: center;
    white-space: normal;
    word-break: break-all;
    word-wrap: break-word;
  }
.platform-tag {
  margin-top: 4px;
}
  .platform-progress {
    margin-top: 4px;
    font-size: 10px;
    font-weight: 200;
  }
}
</style>
