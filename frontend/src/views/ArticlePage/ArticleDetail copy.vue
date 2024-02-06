<template>
  <t-dialog 
    width="80%"
    :header="articleDetailStore.article.title" 
    v-model:visible="articleDetailStore.visible"
    >
    <div style="height: 400px;">
      <t-table 
        :columns="columns" 
        :data="data" 
        :filter-value="filterValue" 
        row-key="PlatformName"
        @filter-change="onFilterChange" lazy-load>
        <template #Status="{row}">
          <!-- {{ row }} -->
          <t-tag :theme="statusNameListMap[row.Status].theme" variant="light">
          {{statusNameListMap[row.Status].label}}
          </t-tag>
        </template>
      </t-table>
    </div>
  </t-dialog>
</template>

<script setup>
import { Textarea } from "tdesign-vue-next";
import { computed, onMounted, ref, watch } from "vue";
import {useArticleDetailStore} from "@/src/store/article"
const articleDetailStore = useArticleDetailStore()


const statusNameListMap = {
  "上传成功": {  theme: 'success', label: "上传成功"},
  "上传失败": { theme: 'danger', label: "上传失败" },
  "等待中": { theme: 'default', label: "等待中" },
  // "": { theme: 'default', label: "等待中" },
};
const filterValue = ref({});

const columns = ref([
  {
    colKey: 'PlatformName', title: '平台名称', width:150,
    filter: {
      type: 'single',
      list: computed(() => {
        const tmp = []
        articleDetailStore.imageProgressList.map(item => tmp.indexOf(item.PlatformName) !== -1 ? null : tmp.push(item.PlatformName))
        return tmp
      }),
    }
  },
  { colKey: 'ImageName', title: '图片名称', width: 240,ellipsis: true,},
  { colKey: 'UploadURL', title: '上传URL', ellipsis: true,},
  { 
    colKey: 'Status', title: '上传状态', width: 120,
    filter: {
      type: 'multiple',
      list: computed(() => {
        const tmp = [{label:'全选', checkAll: true}]
        Object.keys(statusNameListMap).forEach(key => {
          tmp.push({label:key, value:key})
        }) 
        return tmp
      }),
      // showConfirmAndReset: true,
    }
  },
])

const data = ref([...articleDetailStore.imageProgressList])

watch(
  ()=> articleDetailStore.imageProgressList,
  (newValue, oldValue) => {
    data.value = [...articleDetailStore.imageProgressList]
  },
  {deep: true}
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
