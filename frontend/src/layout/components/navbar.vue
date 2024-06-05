<script setup lang="ts">
// import Search from "./search/index.vue";
// import Notice from "./notice/index.vue";
import mixNav from "./sidebar/mixNav.vue";
import { useNav } from "@/layout/hooks/useNav";
import Breadcrumb from "./sidebar/breadCrumb.vue";
import topCollapse from "./sidebar/topCollapse.vue";
// import LogoutCircleRLine from "@iconify-icons/ri/logout-circle-r-line";
import Setting from "@iconify-icons/ri/settings-3-line";
import { onMounted, reactive, ref } from "vue";
import MarkdownIt from "markdown-it";
import {
  CheckUpdate,
  GetloginStatus, GetUpdateUpgrader,
  GetUserInfo,
  GetVersion,
  ReadConfig
} from "../../../wailsjs/go/main/App";
import router from "@/router";
import { Monitor, Start, Stop } from "../../../wailsjs/go/main/Program";
import { ElNotification, ElMessageBox } from "element-plus";

const {
  layout,
  device,
  logout,
  onPanel,
  pureApp,
  avatarsStyle,
  toggleSideBar
} = useNav();
const data = reactive({
  dialogVisible: false,
  version: "v1.0.0",
  updateloading: false,
  updateinfo: {
    Code: 0,
    Msg: "",
    Content: ""
  },
  islogin: false,
  username: "",
  avatars: "https://avatars.githubusercontent.com/u/44761321?v=4",
  isrunning: false,
  form: {
    RoomId: 3,
    DanmuLen: 20,
    EntryMsg: "花花机器人进入直播间",
    PKNotice: true,
    InteractWord: true,
    EntryEffect: true,
    ThanksGift: true,
    WelcomeSwitch: true,
    WelcomeString: {
      "123456": "欢迎宇宙无敌最帅的xxx进入直播间"
    },
    WelcomeDanmu: ["欢迎 {user} ~", "欢迎 {user} 木嘛~", "欢迎 {user} 好诶~"],
    RobotName: "花花",
    TalkRobotCmd: "花花",
    RobotMode: "QingYunKe",
    ChatGPT: {
      APIToken: ""
    },
    FocusDanmu: [
      "啾咪~",
      "喜欢可以领牌牌哦~",
      "么么哒~",
      "入股不亏哦~",
      "贴贴~"
    ],
    CronDanmu: false,
    CronDanmuList: [
      {
        Cron: "*/2 * * * *",
        Random: true,
        Danmu: [
          "喜欢主播请关注, 主播带你去致富~",
          "万水千山总是情, 上个舰长行不行~",
          "喜欢主播的小伙伴可以动动小手点个关注~",
          "喜欢主播的小伙伴，点点关注不迷路~",
          "你已经是成熟的观众了，该学会自己上船了~",
          "小礼物和弹幕都是对主播的支持哦，比心心~",
          "有一种关心叫关注，有一种惦记叫入粉",
          "有一种陪伴叫: 加入大航海~",
          "iOS端可关注公众号哗哩哗哩直播姬充值~",
          "万水千山总是情，点个关注行不行~"
        ]
      }
    ]
  }
});
function updateit() {
  data.updateloading = true;
  GetUpdateUpgrader().then(res => {
    ElNotification({
      title: "更新器下载中",
      message: "更新器下载中",
      type: "success"
    });
  });
}
function checkupdates(isbutton?) {
  CheckUpdate().then(res => {
    data.updateinfo = res;
    if (isbutton == true) {
      if (data.updateinfo.Code == 2) {
        ElNotification({
          title: "检查更新失败",
          message: data.updateinfo.Msg,
          type: "error"
        });
      } else {
        data.dialogVisible = true;
        ElNotification({
          title: "检查更新成功",
          message: data.updateinfo.Msg,
          type: "success"
        });
      }
    }

  });
}
const markdown = new MarkdownIt();
onMounted(() => {
  getuserinfo();
  GetVersion().then(res => {
    data.version = res;
  });
  setInterval(() => {  
    Monitor().then(res => {
      data.isrunning = res;
    });
    checkupdates();
    // console.log(data.isrunning)
  }, 5000);
  setInterval(() => {
    ReadConfig().then(res => {
      if (!res.Code) {
        // console.log(res)
        ElNotification({
          title: "读取配置文件失败",
          message: res.Msg,
          type: "error"
        });
      } else {
        data.form = res.Form;
        // console.log(data.form)
      }
    });
  }, 1000);
});

async function getuserinfo() {
  console.log(data);

  await GetloginStatus().then(res => {
    data.islogin = res;
  });
  if (data.islogin) {
    // window.localStorage.setItem("userInfo","true")
    // router.push("/login")
    await GetUserInfo().then(res => {

      console.log('res', res);
      data.avatars = res.Avactor;
      data.username = res.Username;
      if (!data.username) {
        data.username = '无'
      }

      if (data.avatars == 'data:image/gif;base64,') {
        data.avatars = 'https://img1.baidu.com/it/u=2494495472,3629111731&fm=253'
      }
      console.log(data);
    });
  } else {
    window.localStorage.removeItem("userInfo");
    router.push("/login");
    return;
  }
  pgstart();
}
async function pgstart() {
  await Monitor().then(res => {
    data.isrunning = res;
  });
  if (data.isrunning == false) {
    await Start().then(res => {
      data.isrunning = res;
    });
  }
}
async function pgstop() {
  await Monitor().then(res => {
    data.isrunning = res;
  });
  if (data.isrunning == true) {
    await Stop().then(res => {
      data.isrunning = !res;
    });
  }
}
async function restart() {
  await pgstop();
  console.log(data.isrunning);
  if (data.isrunning == false) {
    pgstart();
  }
}
</script>

<template>
  <div class="navbar bg-[#fff] shadow-sm shadow-[rgba(0, 21, 41, 0.08)] dark:shadow-[#0d0d0d]">

    <!-- <topCollapse v-if="device === 'mobile'" class="hamburger-container" :is-active="pureApp.sidebar.opened"
        @toggleClick="toggleSideBar" /> -->

    <!--  <Breadcrumb v-if="layout !== 'mix' && device !== 'mobile'" class="breadcrumb-container" /> -->

    <!-- <mixNav v-if="layout === 'mix'" / -->
    <div class="vertical-header-left">
      <div class="el-dropdown-link">
        <div>
          <img :src="data.avatars" />
        </div>
        <p v-if="data.username" class="dark:text-white">
          {{ data.username }}
        </p>
      </div>
      <div class="el-dropdown-link">
        <div> 当前房间号</div>{{ data.form.RoomId }}
      </div>
      <div class="el-dropdown-link">
        <div> 运行状态</div>
        <el-tag type="success" v-if="data.isrunning">已启动</el-tag>
        <el-tag type="danger" v-else>已停止</el-tag>
      </div>
      <div class="el-dropdown-link">
        <div>当前版本</div>
        <el-tag type="success">{{ data.version }}</el-tag>
      </div>
    </div>
    <div v-if="layout === 'vertical'" class="vertical-header-right">
      <el-button size="small" type="warning" plain @click="restart">重启</el-button>
      <el-button size="small" type="danger" @click="pgstop" plain v-if="data.isrunning">停止</el-button>
      <el-button size="small" type="success" @click="pgstart" plain v-else>启动</el-button>
      <el-button size="small" type="primary" @click="checkupdates(true)" plain>更新</el-button>

      <span class="set-icon navbar-bg-hover" title="打开项目配置" @click="onPanel">
        <IconifyIconOffline :icon="Setting" />
      </span>
    </div>
  </div>
  <el-dialog v-model="data.dialogVisible" title="检查更新" width="50%">
    <span>当前版本：{{ data.version }}</span><br />
    <span>最新版本：{{ data.updateinfo.Code == 1 ? data.updateinfo.Msg : data.version }}</span><br>
    <div v-html="markdown.render(data.updateinfo.Content)"></div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="data.dialogVisible = false">Cancel</el-button>
        <el-button type="primary" :loading=data.updateloading @click="updateit">
          立即更新
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style lang="scss" scoped>
.navbar {
  overflow: hidden;
  z-index: 990;
  height: 95px;
  overflow: hidden;
  box-shadow: 0 0 1px #888;
  border-radius: 10px;
  margin: 10px;
  padding: 10px;
  display: flex;
  justify-content: space-between;

  .hamburger-container {
    float: left;
    height: 100%;
    line-height: 48px;
    cursor: pointer;
  }

  .vertical-header-right {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    min-width: 280px;
    height: 48px;
    color: #000000d9;
  }

  .vertical-header-left {
    display: flex;
    align-content: end;
    flex-wrap: nowrap;
    align-items: flex-end;
    font-size: 14px;
  }

  .breadcrumb-container {
    float: left;
    margin-left: 16px;
    display: flex;
    align-items: flex-end;
    margin-bottom: -10px;
  }
}

.el-dropdown-link {
  display: flex;
  justify-content: center;
  height: 100%;
  padding: 0 10px;
  color: rgba(0, 0, 0, 0.8509803922);
  cursor: pointer;
  flex-wrap: wrap;
  align-content: center;

  div {
    width: 100%;
    display: flex;
    justify-content: center;
  }

  img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    margin-bottom: 6px;
  }
}

.logout {
  max-width: 120px;

  ::v-deep(.el-dropdown-menu__item) {
    display: inline-flex;
    flex-wrap: wrap;
    min-width: 100%;
  }
}
</style>
