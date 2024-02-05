<template>
    <t-dialog  v-model:visible="batchImportStore.visible" header="批量导入" width="60%" @confirm="batchImport">
        <t-form style="width: 100%;">
            <t-row :gutter="[0, 12]">
                <t-col :span="6"> 
                    <t-form-item label="解析格式">
                        <t-select :options="batchImportStore.formatOptions" v-model="batchImportStore.formatSelectValue" />
                    </t-form-item>
                </t-col>
                <t-col :span="6">
                    <t-form-item label="解析至">
                        <t-select :options="batchImportStore.importTypeOptions" v-model="batchImportStore.importType" :disabled="!batchImportStore.typeInsert"/>
                    </t-form-item>
                </t-col>
                <t-col :span="12">
                    <t-textarea v-model:value="batchImportStore.content" :autosize="{ minRows: 10, maxRows: 20 }"
                    :placeholder="batchImportStore.formatPlacehodler[batchImportStore.formatSelectValue]" />
                </t-col>
            </t-row>
        </t-form>
    </t-dialog>
</template>

<script setup>
import { useBatchImportStore } from '@/src/store/platform'
import { MessagePlugin } from "tdesign-vue-next";
const batchImportStore = useBatchImportStore()


function batchImport() {
    // console.log("batchImportArr: ", batchImportArr.value)
    if (batchImportStore.formatSelectValue === "ROWDATA") {
        var lines = batchImportStore.content.split('\n')
        lines.forEach(element => {
            var line = element.split(': ')
            var key = line[0].trim()
            // var value = line.slice(1).join(': ').trim()
            var value = line[1].trim()
            //   batchImportStore.contentArr.push({ Dynamic: false, Key: key, Value: value })
            var tmp = { Key: key, Value: value }
            if (batchImportStore.typeInsert)  {
                tmp['Type'] = batchImportStore.importType
            }
            batchImportStore.contentArr.push(tmp)
        });
        MessagePlugin.success("导入成功")
        batchImportStore.visible = false
    }
    if (batchImportStore.formatSelectValue  === "JSON") {
        try {
            var jsonOBJ = JSON.parse(batchImportStore.content.toString())
            Object.keys(jsonOBJ).forEach((key) => {
                var tmp = { Key: key, Value: String(jsonOBJ[key]) }
                if (batchImportStore.typeInsert)  {
                tmp['Type'] = batchImportStore.importType
            }
                batchImportStore.contentArr.push(tmp)
            })
            MessagePlugin.success("导入成功")
            batchImportStore.visible = false
        } catch (error) {
            MessagePlugin.error("JSON解析错误: " + error)
        }
    }
}

</script>

<style></style>