<script setup lang="ts">
// import { useI18n } from "vue-i18n";
import Motion from "../utils/motion";
import ReQrcode from "@/components/ReQrcode";
import { useUserStoreHook } from "@/store/modules/user";

// const { t } = useI18n();
</script>

<template>
  <Motion class="-mt-2 -mb-2"> <ReQrcode text="扫码登录" /> </Motion>
  <Motion :delay="200">
    <el-divider>
      <p class="text-gray-500 text-xs">哔哩哔哩APP 扫码登录</p>
    </el-divider>
  </Motion>
<!--  <Motion :delay="150">-->
<!--    <el-button-->
<!--      class="w-full mt-4"-->
<!--      @click=""-->
<!--    >-->
<!--&lt;!&ndash;      {{ t("login.back") }}&ndash;&gt;-->
<!--    </el-button>-->
<!--  </Motion>-->
</template>
