import { defineStore } from "pinia";
import { computed, ref } from "vue";

// 文章组件参数
export const usePlatformStore = defineStore('platform', ()=>{
    const platforms = ref([])
    return {platforms}
})

// 文章详情组件参数
export const useInterfacesStore = defineStore('interface', ()=>{
    const visible = ref(false)
    const title = ref("接口配置")
    const platform = ref({
        Name: "",
        ID: undefined,
        Interfaces: []
    })
    const interfaces = ref([])
    // const options = ref([])
    const options = computed(()=>{
        var tmp = [
            {
              label: "文章信息", value: "0", children:
                [
                  { label: "文章内容", value: "0.article" },
                  { label: "文章名称", value: "0.title" },
                ]
            }
          ]
          console.log("interfacesStore1", platform.value.Interfaces)
          console.log("platform", platform.value)
          
          platform.value.Interfaces.forEach(interface_ => {
            // 分析本接口
            if (interface_.ResponseType != null && interface_.ResponseType === "JSON") {
              try {
                var responseJSON = JSON.parse(interface_.ResponseTemple)
                var node = {
                  label: interface_.Name,
                  value: String(interface_.Serial),
                  children: []
                }
                var child = convertData(responseJSON, String(interface_.Serial))
                node.children = child
                tmp.push(node)
              } catch (error) {
                console.log("load options error: ", error)
              }
            }
        
            // 遍历接口的子接口
            if (interface_.Children != [] && interface_.Children != null) {
              interface_.Children.forEach(subInterface => {
                if (subInterface.ResponseType === "JSON") {
                  var responseJSON = JSON.parse(subInterface.ResponseTemple)
                  var node = {
                    label: "["+interface_.Name+"] "+subInterface.Name,
                    value: String(subInterface.Serial),
                    children: []
                  }
                  var child = convertData(responseJSON, String(subInterface.Serial))
                  node.children = child
                  tmp.push(node)
                }
              })
            }
          })
        
          return tmp
    })
    return {visible, title, platform, interfaces, options}
})

// 转化JSON至Cascader数据,数据值value为1.data.src
function convertData(originalData, prefix = null) {
    const result = [];
    // 递归处理数据
    function handleData(data, parent = null, prefix = null) {
      for (let key in data) {
        let tmp = prefix ? prefix + "." + key : key
        if (typeof data[key] === "object") {
          const item = {
            label: key,
            value: tmp,
            children: []
          };
          if (parent) {
            parent.children.push(item);
          } else {
            result.push(item);
          }
          handleData(data[key], item, tmp);
        } else {
          const item = {
            label: key,
            // value: data[key]
            value: tmp,
          };
          if (parent) {
            parent.children.push(item);
          } else {
            result.push(item);
          }
        }
      }
    }
    handleData(originalData, parent = null, prefix = prefix);
    return result;
  }