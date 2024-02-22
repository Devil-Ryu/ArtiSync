<template>
  <div>
    <t-form>
      <t-row>
        <t-col :span="10">
          <t-row :gutter="[10, 24]">
            <t-col :span="4">
              <t-form-item label="日志等级" name="level">
              <t-select multiple :min-collapsed-num="1" :options="TAB_LIST" v-model="filterValue.level"/>
            </t-form-item>
            </t-col>
            <t-col :span="4">
              <t-form-item label="日志标签" name="tag">
                <t-input
                  type="search"
                  placeholder="请输入日志标签"
                  v-model="filterValue.tag"
                />
              </t-form-item>
            </t-col>
            <t-col :span="4">
              <t-form-item label="日志消息" name="message">
                <t-input
                  type="search"
                  placeholder="请输入关键字"
                  v-model="filterValue.message"
                />
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

    <t-table
      height="600"
      :data="data"
      :columns="columns"
      row-key="datetime"
      style="margin-top: 10px;"
    >   
      <template #level="{ row }">
        <t-tag :theme="statusThemeListMap[row.level]" variant="light">{{statusNameListMap[row.level]}} </t-tag>
      </template>
      <template #tag="{ row }">
        <t-tag :theme="statusThemeListMap[row.level]" variant="light">{{row.tag}} </t-tag>
      </template>
      <template #options="{ row }">
        <t-link theme="primary" @click="logDetail(row)">详情</t-link>
      </template>
    </t-table>
  </div>
  <DetailDialog />
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useLogStore, useLogDetailStore, statusNameListMap, statusThemeListMap } from "@/src/store/log"
import DetailDialog from "./components/DetailDialog.vue"
import { BrowseIcon } from "tdesign-icons-vue-next";
const logStore = useLogStore()
const logDetailStore = useLogDetailStore()

const columns = [
  { title: '时间', colKey: 'datetime', ellipsis: true, width: 180, align: 'center' },
  { title: '等级', colKey: 'level', ellipsis: true, width: 80, align: 'center' },
  { title: '标签', colKey: 'tag', ellipsis: true, width: 150, align: 'center'},
  { title: '消息', colKey: 'message', ellipsis: true, },
  { title: '操作', colKey: 'options', align: 'center', fixed: 'right'},
]


// const tabValue = ref("ALL")
const TAB_LIST = [
{ label: "信息", value: "INF" },
{ label: "调试", value: "DEB" },
{ label: "警告", value: "WAR" },
{ label: "错误", value: "ERR" },
// { label: "全部", value: "ALL" },
]


const data = ref([...logStore.history])
const filterValue = ref({
  level: [],
  tag: "",
  message: "",
})

watch(
  ()=> logStore.history,
  (newValue, oldValue) => {
    data.value = [...logStore.history]
  },
  {deep: true}
)

function logDetail(item) {
  logDetailStore.logData = item
  logDetailStore.visible = true
}

function reset() {
  filterValue.value= {
    level: [],
    tag: "",
    message: "",
  }
  data.value = [...logStore.history]
}

function qeury() {
  console.log("filterValue", filterValue.value)
  setTimeout(()=>{
    const newData = logStore.history.filter((item) => {
      let result = true
      console.log("filterValue.value.level.length !== 0", filterValue.value.level.length )
      if (filterValue.value.level.length !== 0) {
        result = filterValue.value.level.indexOf(item.level) !== -1
      }
      if (filterValue.value.tag !== "") {
        result = item.tag.indexOf(filterValue.value.tag) !== -1
      }
      if (filterValue.value.message !== "") {
        result = item.message.indexOf(filterValue.value.message) !== -1
      }
      return result
    })
    data.value = newData
  }, 100)
  
}

</script>

<style>
.secondary-msg-list {
height: 70vh;

.t-list-item {
  cursor: pointer;
  padding: 13px 24px 13px 0;
  /* padding: 0 24px 0 0; */
  /* padding: 0; */

  &:hover {
    background-color: var(--td-bg-color-container-hover);

    .msg-date {
      display: none;
    }

    .msg-action {
      display: flex;
      align-items: center;

      &-icon {
        display: flex;
        align-items: center;
      }
    }
  }

  .t-tag {
    margin-right: var(--td-comp-margin-s);
    margin-left: 0;
  }
}


.msg-action {
  display: none;

  .set-detail-icon {
      margin-right: 12px;
    }
}

&__empty-list {
  min-height: 443px;
  padding-top: 170px;
  text-align: center;
  color: var(--td-text-color-primary);
}
}



</style>