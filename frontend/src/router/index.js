import { createRouter, createWebHistory } from "vue-router";
import ArticlePage from "../views/ArticlePage/ArticlePage.vue"
import PlatformPage from "../views/PlatformPage/PlatformPage.vue"


// vue项目自带路由
const routes = [
  {
    path: "/",
    redirect: "/article"
  },
  {
    path: "/article",
    name: "Article",
    component: ArticlePage,
    meta:{
        keepAlive:true // 需要缓存
    }
  },
  {
    path: "/platform",
    name: "Platform",
    component: PlatformPage
  },{
    path: "/logs",
    name: "Log",
    component: ()=>import("../views/Logs/LogPage.vue"),
    meta:{
      keepAlive:false // 需要缓存
    }
  },
  {
    path: "/system",
    name: "System",
    component: ()=>import("../views/System/System.vue"),
    meta:{
      keepAlive:false // 需要缓存
    }
  }
];

// const routers = [...routes];

const router = createRouter({
  history: createWebHistory(),
  routes: [...routes]
});

export default router;

