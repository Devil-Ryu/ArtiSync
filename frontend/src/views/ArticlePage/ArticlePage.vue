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
          :data="articleStore.tableData"
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
import {GetArticleImageProgress, LoadArticles, OpenDir, Run} from "@/wailsjs/go/api/ATController"
import {EventsOn} from "@/wailsjs/runtime/runtime";
import { MessagePlugin } from "tdesign-vue-next";
import { LoadJSONFile } from "@/wailsjs/go/api/DBController";
import {useArticleStore, useArticleDetailStore} from "@/src/store/article"
import { useConfigStore } from "@/src/store/config";

const configStore = useConfigStore()
const articleStore = useArticleStore()
const articleDetailStore = useArticleDetailStore()

const statusNameListMap = {
  "运行完成": {  theme: 'success', label: "运行完成"},
  "上传成功": {  theme: 'success', label: "上传成功"},
  "上传失败": { theme: 'danger', label: "上传失败" },
  "处理中": { theme: 'primary', label: "处理中" },
  "等待中": { theme: 'default', label: "等待中" },
};

const columns = ref([
  {colKey: 'Index', title: '序号', width: 60, align: 'center'},
  {colKey: 'Title', title: '文章名称', ellipsisTitle: true},
  {colKey: 'Progress', title: '图片上传进度', width: 240,
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
          articleStore.tableData = response.Data
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


// 刷新文章进度
function RefreshArticleProgress(articleIndex, articleStatus) {
  var imageProgressList = articleStore.imageProgressMap[articleIndex]
  var article = articleStore.tableData[articleIndex]
  var imageCount = imageProgressList.length
  var imageComplete = imageProgressList.filter(item => item.Status === "上传成功" ||  item.Status === "上传失败").length
  article.Progress = parseFloat((imageComplete / imageCount * 100).toFixed(2)) 

  // 更新文章状态
  article.Status = articleStatus
  console.log("article", article)

}

// 加载文章详情
function loadArticleDetail(article) {
  // 更新文章详情的Index
  GetArticleImageProgress(article.Index-1).then(response => {
    if (response.StatusCode == 200) {
      // 更新文章字典上传进度
      articleStore.imageProgressMap[article.Index-1] = response.Data
      // 更新文章详情页
      articleDetailStore.article.title = article.Title + "-上传进度"
      // articleDetailStore.imageProgressList = response.Data == null ? []:response.Data
      articleDetailStore.imageProgressList = response.Data 
      // 打开文章详情弹窗
      articleDetailStore.visible = true
    } else {
      MessagePlugin.error(response.Message)
    }
  })
}

EventsOn("UpdateArticleStatus", async (articleIndex, imageProgressList, articleStatus)=>{
  // 更新文章Map
  articleStore.imageProgressMap[articleIndex] = imageProgressList
  // articleStore.imageProgressMap[articleIndex] = imageProgressList == null ? []:imageProgressList
  // 更新文章详情
  articleDetailStore.imageProgressList = imageProgressList
  // 刷新文章进度
  RefreshArticleProgress(articleIndex, articleStatus)
})




</script>

<style scoped>

</style>