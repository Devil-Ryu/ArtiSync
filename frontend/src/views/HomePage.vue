<template>
  <t-layout style="height: 100%;">
    <!-- <t-header>
      <t-head-menu value="item1">
        <template #logo>
          <p style="color: var(--td-brand-color); font-size: 24px; letter-spacing: 10px;">
            ArtiSync
          </p>
        </template>
        <template #operations>
          <a href="javascript:;">
            <SearchIcon class="t-menu__operations-icon" />
          </a>
          <a href="javascript:;">
            <NotificationFilledIcon class="t-menu__operations-icon" />
          </a>
          <a href="javascript:;">
            <HomeIcon class="t-menu__operations-icon" />
          </a>
        </template>
      </t-head-menu>
    </t-header> -->
    <t-layout class="content">
      <t-aside class="content-left">
        <t-menu theme="light" v-model="asideIndex" width="200px" @change="onChange">
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
              <FormIcon />
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
      <t-layout class="content-right">
        <t-content>
          <div class="content-view">
            <router-view></router-view>
          </div>
        </t-content>
        <div class="content-footer">Copyright @ 2019-{{ new Date().getFullYear() }} Devil-Ryu. All Rights Reserved</div>
      </t-layout>
    </t-layout>
  </t-layout>
</template>
<script setup>
import { FormIcon, LettersRIcon, LettersTIcon, LettersIIcon, LettersSIcon, LettersYIcon, LettersNIcon, LettersCIcon, SearchIcon, NotificationFilledIcon, HomeIcon, BookOpenIcon, ControlPlatformIcon, NotificationIcon, SettingIcon } from "tdesign-icons-vue-next";
import { GetConfigFilePath, LoadJSONFile } from "@/wailsjs/go/api/DBController.js";
import { ref, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { useLogStore } from "@/src/store/log";
import { EventsOn } from "@/wailsjs/runtime/runtime";
import { useConfigStore } from "@/src/store/config"
import { useInterfaceRecordsStore } from "@/src/store/platform";

const configStore = useConfigStore()
const interfaceRecordsStore = useInterfaceRecordsStore()
const asideIndex = ref('article')
const logger = useLogStore()

EventsOn("PushLog", (result) => {
  logger.history.push(result)
})

onMounted(() => {
  loadSystemConfig()
})

function onChange(menuValue) {
  console.log("当前页面[menu]: ", menuValue)
  console.log("当前页面[menu]: ", menuValue==="records")
  if (menuValue === "records") {
    interfaceRecordsStore.setFilters({date_time: []})
  }
}

// 加载系统设置
function loadSystemConfig() {
  GetConfigFilePath().then((configFilePath) => {
    configStore.sysPath = configFilePath
    LoadJSONFile(configFilePath).then((result) => {
      configStore.systemConfig = result
      console.log("加载系统配置: ", result)
    }).catch((err) => {
      MessagePlugin.error(err)
    })
  }).catch((err) => {
    MessagePlugin.error(err)
  })

}
</script>


<style>

.content {

  .content-left {
    width: 200px;
    border-top: 1px solid var(--td-component-border);
    background-color: aqua;
  }

  .content-right {
    margin: 1px;
    padding: 10px;
    height: inherit;
    /* background-color: white; */

    .content-view {
      height: 100%;
      margin: 0;
    }

    .content-footer {
      display: flex;
      margin: 0;
      justify-content: center;
    }
  }


}
</style>