import { defineStore } from "pinia";
import { computed, ref } from "vue";

// 文章组件参数
export const useArticleStore = defineStore('article', () => {
    const articleList = ref([])  // 文章列表
    const articleDetailVisible = ref(false)  // 文章详情是否可见
    const articleDetailIndex = ref(0)

    // 文章详情基本信息
    const articleDetailBasicInfo = computed(()=>{
        if (articleList.value[articleDetailIndex.value] !==  undefined &&  articleList.value[articleDetailIndex.value].BasicInfo !== undefined) {
            // console.log("articleList.value", articleList.value)
            articleList.value[articleDetailIndex.value].BasicInfo.Status = articleList.value[articleDetailIndex.value].Status
            return articleList.value[articleDetailIndex.value].BasicInfo
        } else {
            return {
                Title: "未知",
                Type: "未知",
                ImageNum: "未知",
                InterfaceNum: "未知",
                InterfaceRecordID: "待运行后生成",
                Status: "待发布",
                Platforms: []
            }
        }

    })

    // 文章详平台信息（字典，后端同步到前端）
    const articleDetailPlatformsInfo =  computed(()=>{
        if (articleList.value[articleDetailIndex.value] !==  undefined &&  articleList.value[articleDetailIndex.value].PlatformsInfo !== undefined) {
            return articleList.value[articleDetailIndex.value].PlatformsInfo
        } else {
            return {}
        }
    })
    // 文章详平台信息（列表）
    const articleDetailPlatformsList = computed(() => {
        var result = []
        Object.keys(articleDetailPlatformsInfo.value).forEach(key => {
            var tmp = articleDetailPlatformsInfo.value[key]
            tmp.Name = key
            result.push(tmp)
        })

        // 根据Index进行排序
        result.sort((a, b) => parseInt(a.Index) - parseInt(b.Index))

        return result
    })

    return {
        articleList,
        articleDetailVisible,
        articleDetailIndex,
        articleDetailBasicInfo,
        articleDetailPlatformsInfo,
        articleDetailPlatformsList,
    }
})

export const statusNameListMap = {
    "发布中": { theme: 'primary', label: "发布中" },
    "待发布": { theme: 'default', label: "待发布" },
    "发布成功": { theme: 'success', label: "发布成功" },
    "发布失败": { theme: 'danger', label: "发布失败" },
    "运行成功": { theme: 'success', label: "运行成功" },
    "运行失败": { theme: 'danger', label: "运行失败" },
    "等待中": { theme: 'default', label: "等待中" },
    "运行中": { theme: 'primary', label: "运行中" },
}; 