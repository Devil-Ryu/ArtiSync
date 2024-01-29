import { defineStore } from "pinia";
import { ref } from "vue";

// 文章组件参数
export const useArticleStore = defineStore('article', ()=>{
    const tableData = ref([])  // 表格数据
    const imageProgressMap = ref({})  // 图片上传进度字典
 
    // 文章上传进度 1.上传图片 2.上传文章
    // 文章基本信息
    // 文章名称，图片数量，上传平台数量，平台名称
    // 

    return {tableData, imageProgressMap}
})

// 文章详情组件参数
export const useArticleDetailStore = defineStore('articleDetail', ()=>{
    const article = ref({
        title: "文章详情"
    })
    const visible = ref(false)
    const imageProgressList = ref([])
    return {visible, article, imageProgressList}
})