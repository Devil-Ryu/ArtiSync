<template>
  <t-drawer size="100%" :header="false" :footer="true" placement="left" :confirm-btn="null" cancel-btn="返回" v-model:visible="articleDetailStore.visible">
    <div  style="height: 400px; width: 100%; ">
      <t-card>
        <template #title>
          <div style="margin-left: 14px; ">基本信息</div>
        </template>
        <t-descriptions colon style="margin-top: -20px;" size="small" :label-style="{ width: '1%' }"
          :content-style="{ width: '5%' }" :column="2">

          <t-descriptions-item label="文章名称">测试用文章</t-descriptions-item>
          <t-descriptions-item label="文章类型"><t-tag theme="warning">Markdown</t-tag></t-descriptions-item>

          <t-descriptions-item label="图片数量">7</t-descriptions-item>
          <t-descriptions-item label="接口数量">5</t-descriptions-item>

          <t-descriptions-item label="上传状态"><t-tag theme="primary">上传中</t-tag></t-descriptions-item>
          <t-descriptions-item label="上传平台">知乎、CSDN、博客园</t-descriptions-item>
        </t-descriptions>
      </t-card>

      <!-- 平台上传进度 -->
      <t-card style="margin-top: 10px;">
        <template #title>
          <div style="margin-left: 14px; ">平台上传进度</div>
        </template>
        <t-row :gutter="[0, 24]">
          <t-col :span="colSpan">
            <t-progress theme="circle" percentage="80">
              <template #label>
                <div>知乎</div>
                <!-- <div style="font-size: 12px;font-weight: 300; text-align: center; vertical-align:middle;">上传中: 80%</div> -->
                <t-link theme="primary">查看详情</t-link>
              </template>
            </t-progress>
          </t-col>
          <t-col :span="colSpan">
            <t-progress theme="circle" percentage="0">
              <template #label>
                <div>CSDN</div>
                <!-- <div><t-tag theme="primary" variant="light">上传完成-80%</t-tag></div> -->
                <t-link theme="primary">查看详情</t-link>
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
                    <t-select multiple :min-collapsed-num="1" :options="TAB_LIST" v-model="filterValue.level" />
                  </t-form-item>
                </t-col>
                <t-col :span="4">
                  <t-form-item label="接口名称" name="interface">
                    <t-input type="search" placeholder="请输入接口名称" v-model="filterValue.tag" />
                  </t-form-item>
                </t-col>
                <t-col :span="4">
                  <t-form-item label="运行状态" name="status">
                    <t-select multiple :min-collapsed-num="1" :options="TAB_LIST" v-model="filterValue.level" />
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
            <t-tag :theme="statusNameListMap[row.Status].theme" variant="light">
              {{ statusNameListMap[row.Status].label }}
            </t-tag>
          </template>
        </t-table>
      </t-card>
    
    </div>


  </t-drawer>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useArticleDetailStore } from "@/src/store/article"
const articleDetailStore = useArticleDetailStore()

const colSpan = ref(2)


const statusNameListMap = {
  "上传成功": { theme: 'success', label: "上传成功" },
  "上传失败": { theme: 'danger', label: "上传失败" },
  "等待中": { theme: 'default', label: "等待中" },
  // "": { theme: 'default', label: "等待中" },
};
const filterValue = ref({});

const columns = ref([
  { colKey: 'PlatformName', title: '平台名称', width: 150, },
  { colKey: 'InterfaceName', title: '接口名称', ellipsis: true, },
  { colKey: 'InterfaceURL', title: '接口URL', width: 240, ellipsis: true, },
  { colKey: 'Status', title: '运行状态', ellipsis: true, },
])

const data = ref([...articleDetailStore.imageProgressList])

watch(
  () => articleDetailStore.imageProgressList,
  (newValue, oldValue) => {
    data.value = [...articleDetailStore.imageProgressList]
  },
  { deep: true }
)


const request = (filters) => {
  const timer = setTimeout(() => {
    clearTimeout(timer);
    const newData = articleDetailStore.imageProgressList.filter((item) => {
      let result = true;
      if (result && filters.PlatformName) {
        result = item.PlatformName === filters.PlatformName;
      }
      if (result && filters.Status) {
        result = filters.Status.indexOf(item.Status) !== -1
      }
      return result;
    });
    data.value = newData;
  }, 100);
};

const onFilterChange = (filters, ctx) => {
  console.log('filter-change', filters, ctx);
  filterValue.value = {
    ...filters,
  };
  console.log(filters);
  request(filters);
};

</script>
