import { defineStore } from "pinia";
import { ref } from "vue";

// 文章组件参数
export const useConfigStore = defineStore('config', ()=>{
    // 默认配置
    const defaultConfig = {
        dbPath: "backend/models/test.db",
        imagePath: "",
        imageSelect: "相对文章目录",
        proxyURL: "http://127.0.0.1:8080",
        requestSleep: 1
    }

    const systemConfig = ref({
        dbPath: "backend/models/test.db",
        imagePath: "",
        imageSelect: "相对文章目录",
        proxyURL: "http://127.0.0.1:8080",
        requestSleep: 1
    })
    
    const sysPath = ref("正在获取中")
    return {defaultConfig, systemConfig, sysPath}
})

