<template>
  <div>
    <t-row justify="end">
      <t-space :size="10"  style="margin:10px;">
        <t-button variant="outline" theme="primary" @click="loadData">导入数据</t-button>
        <t-button theme="primary" @click="Run">一键发布</t-button>
      </t-space>
    </t-row>
    <div style="height: 500px;">
      <t-table
          :columns="columns"
          :data="articleStore.articleList"
          row-key="Index">
          <template #operation="{row}">
            <t-button  
              variant="text" 
              theme="primary" 
              size="small" 
              @click="loadArticleDetail(row)">
              详情
            </t-button>
          </template>
      </t-table>
    </div>
    
    <ArticleDetail />
  </div>
  
</template>

<script setup lang="jsx">
import {onMounted, ref} from "vue";
import ArticleDetail from "./ArticleDetail.vue";
import { LoadArticles, OpenDir, Run} from "@/wailsjs/go/api/ATController"
import {EventsOn} from "@/wailsjs/runtime/runtime";
import { MessagePlugin } from "tdesign-vue-next";
import {useArticleStore, statusNameListMap} from "@/src/store/article"
import { useConfigStore } from "@/src/store/config";
import { useInterfaceRecordsStore } from "@/src/store/platform";
import RecordDetailDialog from "../Components/RecordDetailDialog/RecordDetailDialog.vue";
const interfaceRecordsStore = useInterfaceRecordsStore()
const configStore = useConfigStore()
const articleStore = useArticleStore()

// const statusNameListMap = {
//   "运行完成": {  theme: 'success', label: "运行完成"},
//   "上传成功": {  theme: 'success', label: "上传成功"},
//   "上传失败": { theme: 'danger', label: "上传失败" },
//   "处理中": { theme: 'primary', label: "处理中" },
//   "等待中": { theme: 'default', label: "等待中" },
// };

const columns = ref([
  {colKey: 'Index', title: '序号', width: 60, align: 'center'},
  {colKey: 'Title', title: '文章名称', ellipsisTitle: true},
  {colKey: 'Progress', title: '上传进度', width: 240,
    cell: (h, { row}) => {
    return (<t-progress theme="line" percentage={row.Progress} />)
  }},
  {colKey: 'Status', title: '状态', width: 80, align: 'center',
    cell: (h, { row }) => {
    return (<t-tag theme={statusNameListMap[row.Status].theme} variant="light">{statusNameListMap[row.Status].label}</t-tag>)
    },
  },
  {colKey: 'operation', title: '操作', align: 'center', width: 80},
])



// 函数区
function loadData() {
  // 选择文章
  OpenDir().then((response) => {
    if (response.StatusCode == 200) {
      // MessagePlugin.info("文章目录: "+response.Data)
      var selectedDir = response.Data

      // 如果没选择文件夹则返回
      if (selectedDir === "") {
        return
      }

      // 读取系统配置
      var config = configStore.systemConfig
      var imagePath = ""
      if (config.imageSelect == "相对文章目录") {
        imagePath = selectedDir + config.imagePath
      }
      if (config.imageSelect == "固定图片目录") {
        imagePath = config.imagePath
      }
      // MessagePlugin.info("图片目录: "+imagePath)
      LoadArticles(selectedDir, imagePath).then((response) => {
        if (response.StatusCode == 200) {
          for (let i = 0; i < response.Data.length; i++) {
            response.Data[i]["Index"] = i+1
          }
          articleStore.articleList = response.Data
        } else {
          MessagePlugin.error("加载文章错误: "+response.Message)
        }
      }).catch((err)=>{
        MessagePlugin.error("加载配置错误: "+err)
      })
      
    } else {
      MessagePlugin.error("打开目录错误: ", response.Message)
    }
    
  })
}


// 加载文章详情
function loadArticleDetail(article) {
  // 获取文章详情

  console.log("article", article)
  articleStore.articleDetailIndex = article.Index - 1
  articleStore.articleDetailBasicInfo = article.BasicInfo
  articleStore.articleDetailBasicInfo.Status = article.Status
  articleStore.articleDetailPlatformsInfo = article.PlatformsInfo

  console.log("articleDetailBasicInfo", articleStore.articleDetailBasicInfo)

  articleStore.articleDetailVisible = true
}


EventsOn("UpdateArticleDetail", async (articleIndex, BasicInfo, PlatformsInfo, Status, Progress)=>{
  // 更新文章列表中的状态
  articleStore.articleList[articleIndex].BasicInfo = BasicInfo  // 更新基础信息
  articleStore.articleList[articleIndex].PlatformsInfo = PlatformsInfo  // 更新平台进度信息
  articleStore.articleList[articleIndex].Status = Status  // 更新文章状态
  articleStore.articleList[articleIndex].Progress = parseFloat(Progress).toFixed(2)  // 更新文章进度

  // 同步更新文章详情中的状态
  articleStore.articleDetailBasicInfo = BasicInfo
  articleStore.articleDetailBasicInfo.Status  = Status
  articleStore.articleDetailPlatformsInfo = PlatformsInfo
})




</script>

<style scoped>

</style>