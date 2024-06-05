<script setup lang="ts">
import { onBeforeMount, reactive } from "vue";

import { ref } from "vue";
import type { FormInstance, FormRules, TabsPaneContext } from "element-plus";
import { ReadConfig, WriteConfig } from "../../../wailsjs/go/main/App";
import { ElNotification } from "element-plus";
import { Monitor, Start, Stop } from "../../../wailsjs/go/main/Program";
import { Menu, Setting } from '@element-plus/icons-vue'
import { onActivated } from "vue";
import { ElMessage } from 'element-plus'
defineOptions({
  name: "Welcome"
});
const formRef = ref<FormInstance>();
const data = reactive({
  isrunning: false,
  savestatus: false,
  savemsg: "",
  tabledata: [],
  KeywordReplyTemp: [],
  keywordDialogData: {
    keyword: "",
    replyword: ""
  },
  keywordReplyDialogtitle: "",
  keywordReplyDialogVisible: false,
  keywordReplyDialogAction: false,
  dialogdata: {
    uid: "",
    msg: ""
  },
  dialogtitle: "",
  dialogVisible: false,
  dialogAction: false,
  form: {
    RoomId: 3,
    DanmuLen: 20,
    EntryMsg: "花花机器人进入直播间",
    PKNotice: true,
    BlindBoxProfitLossStat: true,
    InteractWord: true,
    EntryEffect: true,
    GoodbyeInfo: true,
    InteractSelf: true,
    ThanksShare: true,
    ThanksGift: true,
    ThanksFocus: true,
    ThanksGiftTimeout: 3,
    WelcomeSwitch: true,
    WelcomeString: {
      "123456": "欢迎宇宙无敌最帅的xxx进入直播间"
    },
    WelcomeDanmu: ["欢迎 {user} ~", "欢迎 {user} 木嘛~", "欢迎 {user} 好诶~"],
    RobotName: "花花",
    TalkRobotCmd: "花花",
    RobotMode: "QingYunKe",
    FuzzyMatchCmd: false,
    ChatGPT: {
      APIToken: "",
      APIUrl: "",
      Prompt: "",
      Limit: false
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
    ],
    WelcomeBlacklistWide: [],
    WelcomeBlacklist: [],
    InteractWordByTime: false,
    WelcomeDanmuByTime: [],
    SignInEnable: true,
    DBPath: "./db",
    DBName: "sqliteDataBase.db",
    DrawByLot: true,
    ShowBlockMsg: true,
    KeywordReply: false,
    KeywordReplyList: {}
    // ForeignLanguageTranslationInChinese: {
    //   Enabled: false,
    //   AppID: "",
    //   SecretKey: ""
    // }
  }
});

async function saveConfig() {
  if(!data.form.RoomId || data.form.RoomId==3 ){
    ElMessage.warning('直播间号错误')
    return;
  }
  if (data.dialogVisible) {
    data.dialogVisible = false
  }
  data.savestatus = false;
  data.savemsg = "";
  console.log(data.form.ChatGPT.APIUrl);
  if (data.form.ChatGPT.APIUrl.length == 0) {
    data.form.ChatGPT.APIUrl = "https://api.openai.com/v1";
  }
  await WriteConfig(JSON.stringify(data.form)).then(res => {
    data.savestatus = res.Code;
    data.savemsg = res.Msg;
  });
  if (data.savestatus) {
    ElNotification({
      title: "操作成功",
      message: "操作成功，正在重启机器人",
      type: "success"
    });
    restart();
  } else {
    ElNotification({
      title: "保存失败",
      message: "请修改配置直到保存正确为止，错误信息：" + data.savemsg,
      type: "error"
    });
  }
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

const rules = reactive<FormRules>({
  RoomId: [
    { required: true, message: "此项必填" },
    { type: "number", message: "此处必须为数字" }
  ],
  DanmuLen: [
    { required: true, message: "此项必填" },
    { type: "number", message: "此处必须为数字" }
  ]
});
function stringChange(str: any, num: any) {
  if (str) {
    if (str.length > num) {
      return str.slice(0, num) + '...'
    } else {
      str.slice
    }
  } else {
    return ''
  }
}
console.log('erer5345');
onActivated(() => {
  ReadConfig().then(res => {
    if (!res.Code) {
      console.log(res);
      ElNotification({
        title: "读取配置文件失败",
        message: res.Msg,
        type: "error"
      });
    } else {
      data.form = res.Form;
      if (!data.form.RoomId || data.form.RoomId == 3) {
        data.dialogVisible = true;
      }
      if (data.form.ChatGPT.APIUrl.length == 0) {
        data.form.ChatGPT.APIUrl = "https://api.openai.com/v1";
      }
    }
  });
});

</script>

<template>
  <el-scrollbar class="my-scrollbar">
    <div class="card-content">
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />基础信息
          </span>
          <router-link to="/base">
            <Setting class="icon2" />
          </router-link>
        </template>
        <p>绑定直播间号<span>{{ data.form.RoomId }}</span></p>
        <p>弹幕最大长度<span>{{ data.form.DanmuLen }}</span></p>
        <p>花花进入消息<span>{{ stringChange(data.form.EntryMsg, 6) }}</span></p>
        <p>礼物感谢频率<span>{{ data.form.ThanksGiftTimeout + '秒' }}</span></p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />弹幕类别开关
          </span>
        </template>
        <p>特效入场欢迎<el-switch v-model="data.form.EntryEffect" size="small" @click="saveConfig" /></p>
        <p>礼物感谢<el-switch v-model="data.form.ThanksGift" size="small" @click="saveConfig" /></p>
        <p>PK提醒<el-switch v-model="data.form.PKNotice" size="small" @click="saveConfig" /></p>
        <p>禁言提醒<el-switch v-model="data.form.ShowBlockMsg" size="small" @click="saveConfig" /></p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />新增弹幕开关
          </span>
        </template>
        <p>下播感谢<el-switch v-model="data.form.GoodbyeInfo" size="small" @click="saveConfig" /></p>
        <p>拉黑自己<el-switch v-model="data.form.InteractSelf" size="small" @click="saveConfig" /></p>
        <p>分享感谢<el-switch v-model="data.form.ThanksShare" size="small" @click="saveConfig" /></p>
        <p>盲盒盈亏<el-switch v-model="data.form.BlindBoxProfitLossStat" size="small" @click="saveConfig" /></p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />欢迎弹幕
          </span>
          <div style="display: flex;">
            <el-switch v-model="data.form.InteractWord" size="small" @click="saveConfig" />
            <router-link to="/welcomebarrage">
              <Setting class="icon2" />
            </router-link>
          </div>
        </template>
        <p>分时间段欢迎
          <span style="display: flex;">
            <el-switch v-model="data.form.InteractWordByTime" size="small" @click="saveConfig" />
            <router-link to="/welcomebarrage">
              <Setting class="icon2" />
            </router-link>
          </span>
        </p>
        <p>随机欢迎弹幕
          <router-link to="/welcomebarrage">
            <Setting class="icon2" />
          </router-link>
        </p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />关键词回复
          </span>
          <div style="display: flex;">
            <el-switch v-model="data.form.KeywordReply" size="small" @click="saveConfig" />
            <router-link to="/keywordreply">
              <Setting class="icon2" />
            </router-link>
          </div>
        </template>
        <p>关键词<span>{{
          data.KeywordReplyTemp[0] ? data.KeywordReplyTemp[0].key.length > 6 ? data.KeywordReplyTemp[0].key.slice(0,
            6) + '...' : data.KeywordReplyTemp[0].key : '' }}等</span>
        </p>
        <p>回复内容<span>{{ data.KeywordReplyTemp[0] ? data.KeywordReplyTemp[0].value.length > 6 ?
          data.KeywordReplyTemp[0].value.slice(0, 6) + '...' : data.KeywordReplyTemp[0].value : '' }}等</span></p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />指定人欢迎
          </span>
          <div style="display: flex;">
            <el-switch v-model="data.form.WelcomeSwitch" size="small" @click="saveConfig" />
            <router-link to="/assigner">
              <Setting class="icon2" />
            </router-link>
          </div>
        </template>
        <p>用户uid<span>{{
          data.tabledata[0] ? data.tabledata[0].key.length > 6 ? data.tabledata[0].key.slice(0,
            6) + '...' : data.tabledata[0].key : '' }}等</span>
        </p>
        <p>欢迎语<span>{{ data.tabledata[0] ? data.tabledata[0].value.length > 6 ?
          data.tabledata[0].value.slice(0, 6) + '...' : data.tabledata[0].value : '' }}等</span></p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />ai机器人
          </span>
          <div style="display: flex;">
            <router-link to="/robots">
              <Setting class="icon2" />
            </router-link>
          </div>
        </template>
        <p>机器人名称<span>{{ data.form.RobotName }}</span></p>
        <p>触发关键字<span>{{ data.form.TalkRobotCmd }}</span></p>
        <p>服务提供商<span>{{ data.form.RobotMode }}</span></p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />关注答谢语
          </span>
          <div style="display: flex;">
            <el-switch v-model="data.form.ThanksFocus" size="small" @click="saveConfig" />
            <router-link to="/thank">
              <Setting class="icon2" />
            </router-link>
          </div>
        </template>
        <p v-for="(items, index) in data.form.FocusDanmu.slice(0, 3)">
          {{ index + 1 + '. ' + data.form.FocusDanmu[index] }}
        </p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />定时弹幕
          </span>
          <div style="display: flex;">
            <el-switch v-model="data.form.CronDanmu" size="small" @click="saveConfig" />
            <router-link to="/timedbarrage">
              <Setting class="icon2" />
            </router-link>
          </div>
        </template>
        <p v-for="(items, index) in data.form.CronDanmuList[0].Danmu.slice(0, 3)">
          {{ index + 1 + '. ' + stringChange(items, 12) }}
        </p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />欢迎黑名单
          </span>
          <router-link to="/blacklist">
            <Setting class="icon2" />
          </router-link>
        </template>
        <p>模糊匹配
          <router-link to="/blacklist">
            <Setting class="icon2" />
          </router-link>
        </p>
        <p>精确匹配
          <router-link to="/blacklist">
            <Setting class="icon2" />
          </router-link>
        </p>
      </el-card>
      <el-card shadow="hover">
        <template #header>
          <span>
            <Menu class="icon" />功能开关
          </span>
        </template>
        <p>抽签<el-switch v-model="data.form.DrawByLot" size="small" @click="saveConfig" /></p>
        <p>签到<el-switch v-model="data.form.SignInEnable" size="small" @click="saveConfig" /></p>
      </el-card>
    </div>
  </el-scrollbar>

  <el-dialog v-model="data.dialogVisible" title="绑定直播间号" width="50%">
    <el-form :model="data.form" ref="formRef" class="el-form" :rules="rules">
      <el-form-item label="直播间号" prop="RoomId">
        <el-input v-model.number="data.form.RoomId" />
      </el-form-item>
    </el-form>
    <template #footer>
      <center>
        <el-button type="success" @click="saveConfig">保存</el-button>
      </center>
    </template>
  </el-dialog>
</template>

<style>
.card-content {
  margin: 20px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;

  /* .el-card:hover{
    box-shadow: 0px 0px 6px 0px rgba(61, 173, 253, 0.555);
  } */
  color: #666;

  .el-card__header {
    padding: 15px;
    font-size: 16px;
    font-weight: 700;
    display: flex;
    justify-content: space-between;

    >span {
      display: flex;
    }
  }

  .el-card__body {
    padding: 15px;
    font-size: 13px;
    background-color: #f5f7fa;
    height: 100%;

    >p {
      height: 33px;
      line-height: 33px;
      display: flex;
      justify-content: space-between;

      >span {
        color: #666;
      }
    }
  }
}

.icon {
  width: 1.5em;
  height: 1.5em;
  color: #409EFF;
  padding-right: 7px;
}

.icon2 {
  width: 1.0em;
  height: 1.0em;
  color: #666;
  height: 24px;
  line-height: 24px;
  margin-left: 5px;
}

.demo-tabs>.el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}

.el-tabs__header {
  position: sticky;
  margin: 0 0 15px;
  top: 0;
  z-index: 2;
}
</style>
