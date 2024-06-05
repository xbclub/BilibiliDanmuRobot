<script setup lang="ts">
import Motion from "./utils/motion";
import { useRouter } from "vue-router";
import { useNav } from "@/layout/hooks/useNav";
import { useLayout } from "@/layout/hooks/useLayout";
import { ref} from "vue";
import { useDataThemeChange } from "@/layout/hooks/useDataThemeChange";
import { bgblue, bgpink, boy, girl, flowerPink, flowerBlue } from "@/utils/static";
import {
  Female,
  Male,
} from '@element-plus/icons-vue'
import QrCode from "@/views/login/components/qrCode.vue";
defineOptions({
  name: "Login"
});
const router = useRouter();
const girlOrBoy = ref(true);

const { initStorage } = useLayout();
initStorage();
const { dataTheme, dataThemeChange } = useDataThemeChange();
dataThemeChange();
function change(val: any) {
  console.log(val);
  girlOrBoy.value = val
}
const { title } = useNav();
</script>

<template>
  <div class="select-none">
    <img :src="bgpink" class="wave" v-if="girlOrBoy" />
    <img :src="bgblue" class="wave" v-else />
    <div class="flex-c absolute right-5 top-3">
      <!-- 主题 -->
      <el-switch v-model="girlOrBoy" style="--el-switch-on-color: #f88597; --el-switch-off-color: #5ed6de" inline-prompt
        :active-icon="Female" :inactive-icon="Male" @change="change" />
    </div>
    <div class="login-container">
      <div v-if="girlOrBoy">
        <img :src="girl" class="wave" />
      </div>
      <div v-else>
        <img :src="boy" class="wave" />
      </div>
      <div class="login-box">
        <div class="login-form">
          <img :src="flowerPink" class="avatar" v-if="girlOrBoy" />
          <img :src="flowerBlue" class="avatar" v-else />
          <Motion>
            <h2 class="outline-none">花花机器人</h2>
          </Motion>
          <qrCode />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@import url("@/style/login.css");
</style>

<style lang="scss" scoped>
:deep(.el-input-group__append, .el-input-group__prepend) {
  padding: 0;
}
</style>
