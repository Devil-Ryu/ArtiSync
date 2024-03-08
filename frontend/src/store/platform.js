import { defineStore } from "pinia";
import { computed, ref } from "vue";

// 文章组件参数
export const usePlatformStore = defineStore('platform', ()=>{
    const platforms = ref([])
    return {platforms}
})

// 接口详情组件参数
export const useInterfacesStore = defineStore('interface', ()=>{
    const visible = ref(false)
    const title = ref("接口配置")
    const platform = ref({
        Name: "",
        ID: undefined,
        Interfaces: [],
        ReplaceMaps: []
    })
    const interfaces = ref([])
    // const options = ref([])
    const options = computed(()=>{
        var tmp = [
            {
              label: "文章信息", value: "ART", children:
                [
                  { label: "文章名称", value: "ART-TITLE" },
                  { label: "文章内容", value: "ART-CONTENT-STR" },
                ]
            },
            {
              label: "图片信息", value: "IMG", children:
                [
                  { label: "图片URL", value: "IMG-URL" },
                  { label: "图片内容[字符串]", value: "IMG-CONTENT-STR" },
                  { label: "图片内容[16进制]", value: "IMG-CONTENT-HEX" },
                ]
            }
          ]
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

    // 平台接口选项
    const SerialOptions = computed(()=>{
      var result = []
      platform.value.Interfaces.forEach(item=>{
        if (item.IsGroup === true) {
          item.Children.forEach(subItem=>{
            result.push({label:subItem.Name, value:subItem.Serial})
          })
        } else {
          result.push({label:item.Name, value:item.Serial})
        }
        
      })
      return result
    })
    return {visible, title, platform, interfaces, options, SerialOptions}
})

// 接口记录参数
export const useInterfaceRecordsStore = defineStore('interfaceRecords', () => {
  const records = ref([])  // 记录列表
  const filters = ref({date_time: []})
  const detailDialogVisible = ref(false)  // 是否显示
  const curRecord = ref({
    ArticleName: "初次渲染模板",
    DateTime: "2024-02-18 17:43:21",
    ID: 795,
    Name: "上传图片",
    PlatformName: "XXXX",
    RecordID: "ix83sLg7",
    RequestMessage: "POST /api/dddd/ccc HTTP/1.1\r\nHost: test.dddd.cc\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 213415\r\nContent-Type…",
    RequestURL: "https://test.dddd.cc/",
    ResponseMessage: "发送请求错误:XXXXXXX\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1",
    Serial: "WI7Unljn",
    Status: "运行失败",
    Tag: "1",
  })
  const testCaches = ref([])  // 测试接口的缓存

  function setFilters(filters) {
    this.filters = filters
    return filters
  }



  return {
      filters,
      records,
      detailDialogVisible,
      curRecord,
      testCaches,
      setFilters,
      
  }
})

// 批量导入组件
export const useBatchImportStore = defineStore('batchImport', ()=>{
  const visible = ref(false)
  const typeInsert = ref(false)
  const importType = ref("Headers")
  const importTypeOptions = ref([
    { label: 'Path', value: 'Path' },
    { label: 'Query', value: 'Query' },
    { label: 'Headers', value: 'Headers' },
    { label: 'Body', value: 'Body' },
  ])
  const formatOptions = ref([
    { label: 'JSON格式', value: 'JSON' },
    { label: 'ROWDATA格式', value: 'ROWDATA' },
    // { label: 'SQLITE文件格式', value: 'SQLITE' },
  ])
  const formatPlacehodler = ref({
    "JSON": "{'key':'value'}",
    "ROWDATA": "key1=value1\nkey2=value2",
    "SQLITE": "Sqlite文件",
  })
  const formatSelectValue = ref("JSON")
  const content = ref("")
  const contentArr = ref([])

  return {visible, typeInsert, importType, importTypeOptions, formatOptions, formatPlacehodler, formatSelectValue, content, contentArr}


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