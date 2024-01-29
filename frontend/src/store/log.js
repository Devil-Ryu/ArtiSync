import { defineStore } from "pinia";
import { computed, ref } from "vue";

// 日志组件参数
export const useLogStore = defineStore('logger', ()=>{
    const history = ref([
        // {"level":"INF", "datetime":"2023-12-12:08:33:33" ,"tag": "文章控制器", "message":"这是一条测试消息1111"},
        // {"level":"INF", "datetime":"2023-12-12:08:33:33" ,"tag": "文章控制器", "message":"这是一条测试消息2222"},
        // {"level":"DEB", "datetime":"2023-12-12:08:33:33" ,"tag": "文章控制器", "message":"这是一条测试消息1111"},
        // {"level":"ERR", "datetime":"2023-12-12:08:33:33" ,"tag": "文章控制器", "message":"这是一条测试消息1111"},
        // {"level":"WAR", "datetime":"2023-12-12:08:33:33" ,"tag": "文章控制器", "message":"这是一条测试消息1111"},
    ])
    
    // const infMsg = computed(()=>{ return history.value.filter((item)=> item.level === "INF") })
    // const debMsg = computed(()=>{ return history.value.filter((item)=> item.level === "DEB") })
    // const warMsg = computed(()=>{ return history.value.filter((item)=> item.level === "WAR") })
    // const errMsg = computed(()=>{ return history.value.filter((item)=> item.level === "ERR") })
    // return {history, infMsg, debMsg, warMsg, errMsg, statusNameListMap}
    return {history}
})

export const useLogDetailStore = defineStore('logDetail', ()=>{
    const visible = ref(false)
    const logData = ref({})

    return {visible, logData}
})

export const statusNameListMap = {
    "INF":  "信息",
    "DEB":  "调试",
    "WAR":  "警告",
    "ERR":  "错误",
}
export const statusThemeListMap = {
    "INF": 'success', 
    "DEB": 'primary', 
    "WAR": 'warning', 
    "ERR": 'danger',  
}