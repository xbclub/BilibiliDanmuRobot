<script setup lang="ts">
import { getTopMenu } from "@/router/utils";
import { useNav } from "@/layout/hooks/useNav";
import { flowerPinkF, flowerBlueF } from "@/utils/static";

const props = defineProps({
  collapse: Boolean
});

const { title } = useNav();
const topPath = "/welcome";
</script>

<template>
  <div class="sidebar-logo-container" :class="{ collapses: props.collapse }">
    <transition name="sidebarLogoFade">
      <router-link v-if="props.collapse" key="props.collapse" :title="title" class="sidebar-logo-link" :to="topPath">
        <img :src="flowerPinkF" alt="logo" />
        <!-- <span class="sidebar-title">{{ title }}</span> -->
      </router-link>
      <router-link v-else key="expand" :title="title" class="sidebar-logo-link" :to="topPath">
        <img :src="flowerPinkF" alt="logo" />
        <span class="sidebar-title">{{ title }}</span>
      </router-link>
    </transition>
  </div>
</template>

<style lang="scss" scoped>
*{
  --el-color-primary: #4f93ea !important;
}
.sidebar-logo-container {
  position: relative;
  width: 100%;
  height: 48px;
  overflow: hidden;

  .sidebar-logo-link {
    display: flex;
    flex-wrap: nowrap;
    align-items: center;
    height: 100%;
    justify-content: center;

    img {
      display: inline-block;
      height: 32px;
      width: 32px;
    }

    .sidebar-title {
      display: inline-block;
      height: 40px;
      margin: 0px 12px;
      overflow: hidden;
      font-size: 18px;
      font-weight: 500;
      line-height: 40px;
      color: $subMenuActiveText;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}
</style>
