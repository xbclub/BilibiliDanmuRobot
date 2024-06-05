<script setup lang="ts">
import { useGlobal } from "@pureadmin/utils";
import backTop from "@/assets/svg/back_top.svg?component";
import { h, computed, Transition, defineComponent } from "vue";
import { usePermissionStoreHook } from "@/store/modules/permission";
import { RouterView } from 'vue-router'

const props = defineProps({
  fixedHeader: Boolean
});

const { $storage, $config } = useGlobal<GlobalPropertiesApi>();

const keepAlive = computed(() => {
  return $config?.KeepAlive;
});

const transitions = computed(() => {
  return route => {
    return route.meta.transition;
  };
});

const hideTabs = computed(() => {
  return $storage?.configure.hideTabs;
});

const layout = computed(() => {
  return $storage?.layout.layout === "vertical";
});

const getSectionStyle = computed(() => {
  return [
    hideTabs.value && layout ? "margin-top: 115px;" : "",
    !hideTabs.value && layout ? "margin-top: 165px;" : "",
    hideTabs.value && !layout.value ? "margin-top:  115px" : "",
    !hideTabs.value && !layout.value ? "margin-top: 165px;" : "",
    props.fixedHeader ? "" : "margin-top: 0;"
  ];
});

const transitionMain = defineComponent({
  render() {
    return h(
      Transition,
      {
        name:
          transitions.value(this.route) &&
            this.route.meta.transition.enterTransition
            ? "pure-classes-transition"
            : (transitions.value(this.route) &&
              this.route.meta.transition.name) ||
            "fade-transform",
        enterActiveClass:
          transitions.value(this.route) &&
          `animate__animated ${this.route.meta.transition.enterTransition}`,
        leaveActiveClass:
          transitions.value(this.route) &&
          `animate__animated ${this.route.meta.transition.leaveTransition}`,
        mode: "out-in",
        appear: true
      },
      {
        default: () => [this.$slots.default()]
      }
    );
  },
  props: {
    route: {
      type: undefined,
      required: true
    }
  }
});
</script>

<template>
  <section :class="[props.fixedHeader ? 'app-main' : 'app-main-nofixed-header']" :style="getSectionStyle">
    <router-view v-slot="{ Component, route }" >
      <transition>
        <keep-alive>
          <component :is="Component" some-prop="a value" />
        </keep-alive>
      </transition>
    </router-view>
  </section>
</template>

<style scoped>
.app-main {
  z-index: 998;
  height: calc(100% - 175px);
  overflow: hidden;
  box-shadow: 0 0 1px #888;
  border-radius: 10px;
  position: fixed;
  background-color: #fff;
  margin: 10px;
}

.app-main-nofixed-header {
  position: relative;
  width: 100%;
  min-height: 100vh;
}

.main-content {
  margin: 24px;
}
</style>
