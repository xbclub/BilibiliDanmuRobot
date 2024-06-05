<template>
  <el-config-provider :locale="currentLocale">
    <keep-alive>
      <router-view />
    </keep-alive>
    <ReDialog />
  </el-config-provider>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { ElConfigProvider } from "element-plus";
import zhCn from "element-plus/lib/locale/lang/zh-cn";
import { ReDialog } from "@/components/ReDialog";
export default defineComponent({
  name: "app",
  components: {
    [ElConfigProvider.name]: ElConfigProvider,
    ReDialog
  },
  computed: {
    currentLocale() {
      return zhCn;
    }
  }
});
</script>
<style>
.between {
  display: flex;
  justify-content: space-between;
  width: 100%;
}
</style>
