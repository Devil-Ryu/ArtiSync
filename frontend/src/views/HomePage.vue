<template>
  <t-layout style="min-height: 640px;">
    <t-header>
      <t-head-menu value="item1" >
        <template #logo>
          <p style="color: var(--td-brand-color); font-size: 24px; letter-spacing: 10px;">
            ArtiSync
          </p>
          <!-- <t-space break-line style="color: var(--td-brand-color); font-size: 25px; letter-spacing: 10px;" :size="0">
            ArtiSync
          </t-space> -->
        </template>
        <!-- <t-menu-item value="item1" :disabled="true"> 发布文章 </t-menu-item> -->
        <template #operations>
          <a href="javascript:;"><SearchIcon class="t-menu__operations-icon"/></a>
          <a href="javascript:;"><NotificationFilledIcon class="t-menu__operations-icon"/></a>
          <a href="javascript:;"><HomeIcon class="t-menu__operations-icon"/></a>
        </template>
      </t-head-menu>
    </t-header>
    <t-layout>
      <t-aside style="border-top: 1px solid var(--td-component-border)" width="200px">
        <t-menu theme="light" v-model="asideIndex" style="margin-right: 50px;" height="550px" width="100%">
          <t-menu-item value="article" to="/article">
            <template #icon>
              <BookOpenIcon />
            </template>
            发布文章
          </t-menu-item>
          <t-menu-item value="platform" to="/platform">
            <template #icon>
              <ControlPlatformIcon />
            </template>
            平台管理
          </t-menu-item>
          <t-menu-item value="records" to="/records">
            <template #icon>
              <ControlPlatformIcon />
            </template>
            接口记录
          </t-menu-item>
          <t-menu-item value="logs" to="/logs">
            <template #icon>
              <NotificationIcon />
            </template>
            日志中心
          </t-menu-item>
          <t-menu-item value="system" to="/system">
            <template #icon>
              <SettingIcon />
            </template>
            系统设置
          </t-menu-item>
        </t-menu>
      </t-aside>
      <t-layout >
        <t-content >
          <div class="content">
            <router-view></router-view>
          </div>
        </t-content>
        <t-footer class="footer">Copyright @ 2019-{{ new Date().getFullYear() }} Devil-Ryu. All Rights Reserved</t-footer>
      </t-layout>
    </t-layout>
  </t-layout>
</template>
<script setup>
import { LettersAIcon, LettersRIcon, LettersTIcon, LettersIIcon, LettersSIcon, LettersYIcon, LettersNIcon, LettersCIcon, SearchIcon, NotificationFilledIcon, HomeIcon, BookOpenIcon, ControlPlatformIcon, NotificationIcon, SettingIcon} from "tdesign-icons-vue-next";
import { GetConfigFilePath, LoadJSONFile} from "@/wailsjs/go/api/DBController.js";
import { ref, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { useLogStore } from "@/src/store/log";
import { EventsOn } from "@/wailsjs/runtime/runtime";
import { useConfigStore } from "@/src/store/config"

const configStore = useConfigStore()
const asideIndex = ref('article')
const logger = useLogStore()

EventsOn("PushLog", (result)=>{
  logger.history.push(result)
})

onMounted(() => {
    loadSystemConfig()
})

// 加载系统设置
function loadSystemConfig() {
  GetConfigFilePath().then((configFilePath)=>{
    configStore.sysPath = configFilePath
    LoadJSONFile(configFilePath).then((result)=>{
        configStore.systemConfig = result
        console.log("loadSystemConfig: ", result)
    }).catch((err)=>{
        MessagePlugin.error(err)
    })
  }).catch((err)=>{
    MessagePlugin.error(err)
  })
    
}


</script>


<style>

.footer {
  display: flex;
  padding: 0;
  height: 1px;
  font-size: 12px;
  justify-content: center;
  align-items: center;
  /* background-color: #2ba471; */
}

.content {
  margin: 1px;
  padding: 10px;
  height: 100%;
  background-color: white;
}

</style>