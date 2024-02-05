<template>
    <div>
        <t-input-adornment>
            <template #prepend>
                <t-button variant="outline" @click="btnclk">
                    <div v-if="!props.item.Disabled" style="color:cornflowerblue;">T</div>
                    <div v-if="props.item.Disabled" style="color:chocolate;">F</div>
                </t-button>
            </template>
            <t-select class="c-select"  :options="['Headers', 'Query'].map((value)=>{return {value, label:value}})" v-model:value="props.item.Type"/>
        </t-input-adornment>
    </div>
</template>

<script setup>
import { MessagePlugin } from "tdesign-vue-next";
const props = defineProps({
    item: {
        type: Object,
        default: () => {
            return {
                Disabled: {
                    type: Boolean,
                    default: false,
                },
                Type: {
                    type: String,
                    default: "",
                },
                Key: {
                    type: String,
                    default: "",
                }
            }
        }
    }

})

function btnclk() {
    props.item.Disabled = !props.item.Disabled
    if (props.item.Disabled) {
        MessagePlugin.error(props.item.Key+" 已禁用")
    } else {

        MessagePlugin.success(props.item.Key+" 已启用")
    }
}

</script>

<style>
.c-select {
    width: 100px;
}
</style>