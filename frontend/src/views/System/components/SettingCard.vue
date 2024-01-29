<template>
<TCard :shadow="true" style="margin-bottom: 10px;">
    <TCollapse expand-icon-placement="right" borderless>
        <TCollapsePanel>
            <template #header>
                <div> {{ props.title }} </div>
            </template>
           <slot name="content">

           </slot>
        </TCollapsePanel>
    </TCollapse>
</TCard>
</template>

<script setup>
const props = defineProps({
    title: {
        type: String,
        default: "新设置"
    }
})
</script>
