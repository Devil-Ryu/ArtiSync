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
          <t-descriptions-item label="文章类型"><t-tag theme="warning" variant="light-outline">{{
            articleStore.articleDetailBasicInfo.Type
          }}</t-tag></t-descriptions-item>

          <t-descriptions-item label="图片数量">{{ articleStore.articleDetailBasicInfo.ImageNum }}</t-descriptions-item>
          <t-descriptions-item label="接口数量">{{ articleStore.articleDetailBasicInfo.InterfaceNum }}</t-descriptions-item>

          <t-descriptions-item label="发布状态"><t-tag variant="light-outline"
              :theme="statusNameListMap[articleStore.articleDetailBasicInfo.Status].theme">{{
                statusNameListMap[articleStore.articleDetailBasicInfo.Status].label
              }}</t-tag></t-descriptions-item>
          <t-descriptions-item label="发布平台">{{ articleStore.articleDetailBasicInfo.Platforms
          }}</t-descriptions-item>
          <t-descriptions-item label="记录ID">{{ articleStore.articleDetailBasicInfo.InterfaceRecordID
          }}</t-descriptions-item>

        </t-descriptions>
      </t-card>

      <!-- 平台上传进度 -->
      <t-card style="margin-top: 10px;">
        <template #title>
          <div style="margin-left: 14px; ">平台上传进度</div>
        </template>
        <t-row :gutter="[0, 24]">
          <t-col :span="colSpan" v-for="platform in articleStore.articleDetailPlatformsList ">
            <t-progress theme="circle" :percentage="platform.RunProgress">
              <!-- <t-progress theme="circle" :percentage="platform.RunProgress" :status="statusNameListMap[platform.RunStatus].theme"> -->
              <template #label>
                <div class="platform-box">
                  <div class="platform-name">{{ platform.Name }}</div>
                  <t-tag class="platform-tag" variant="light-outline"
                    :theme="statusNameListMap[platform.RunStatus].theme">{{ statusNameListMap[platform.RunStatus].label
                    }}</t-tag>
                  <div class="platform-progress"> {{ platform.InterfaceActualRunNum }} / {{
                    platform.InterfacePredictRunNum }} </div>
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
        <InterfaceRecordPage />
      </t-card>
    </div>
    <RecordDetailDialog />
  </t-drawer>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useArticleStore, statusNameListMap } from "@/src/store/article"
import { useInterfaceRecordsStore } from "@/src/store/platform";
import { GetInterfaceRecords, QueryInterfaceRecords } from "@/wailsjs/go/api/DBController";
import RecordDetailDialog from "../Components/RecordDetailDialog/RecordDetailDialog.vue";
import InterfaceRecordPage from "../InterfaceRecords/InterfaceRecordPage.vue";

const articleStore = useArticleStore()
const colSpan = ref(2)

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
