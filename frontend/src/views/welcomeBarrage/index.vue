<script setup lang="ts">
import { onActivated, reactive } from "vue";
import {
  Delete,
} from '@element-plus/icons-vue'
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
    InteractWord: true,
    EntryEffect: true,
    ThanksGift: true,
    ThanksGiftTimeout: 3,
    WelcomeSwitch: true,
    WelcomeString: {
      "123456": "欢迎宇宙无敌最帅的xxx进入直播间"
    },
    WelcomeDanmu: ["欢迎 {user} ~", "欢迎 {user} 木嘛~", "欢迎 {user} 好诶~"],
    WelcomeTier: false,
    WelcomeNew: ["欢迎呀~{at}恭喜你发现了{anchor}"],
    WelcomeBack: ["哈喽呀~{at}等你很久喽~"],
    WelcomeFan: ["哈喽{at}今天有没有更想{anchor}呀~"],
    WelcomeFanMedalLevel: 0,
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
const WelcomeDanmuByTimeTemplate = [
  {
    Enabled: true,
    Key: "earlymorning",
    Random: true,
    Danmu: [
      "欢迎 {user}, 凌晨的问候~",
      "欢迎 {user}, 早安，幸福满满",
      "欢迎 {user}, 祝梦想成真",
      "欢迎 {user}, 望快乐起步",
      "欢迎 {user}, 望奋斗加油",
      "欢迎 {user}, 祝充满力量",
      "欢迎 {user}, 祝精神旺盛",
      "欢迎 {user}, 望谦虚前行",
      "欢迎 {user}, 祝青春无限",
      "欢迎 {user}, 祝前程似锦"
    ]
  },
  {
    Enabled: true,
    Key: "morning",
    Random: true,
    Danmu: [
      "欢迎 {user}, 早安，美好开始",
      "欢迎 {user}, 早安，成功在望",
      "欢迎 {user}, 早安，加油！",
      "欢迎 {user}, 早安，勇往直前",
      "欢迎 {user}, 早安，工作愉快",
      "欢迎 {user}, 早安，眉开眼笑",
      "欢迎 {user}, 早安，携手前行",
      "欢迎 {user}, 早安，让生活更好",
      "欢迎 {user}, 早安，热情满怀"
    ]
  },
  {
    Enabled: true,
    Key: "latemorning",
    Random: true,
    Danmu: [
      "欢迎 {user}, 上午好，奋斗有力",
      "欢迎 {user}, 上午好，朝气蓬勃",
      "欢迎 {user}, 上午好，加油向前",
      "欢迎 {user}, 上午好，祝福无限",
      "欢迎 {user}, 上午好，暖心问候",
      "欢迎 {user}, 上午好，快乐工作",
      "欢迎 {user}, 上午好，信心十足",
      "欢迎 {user}, 上午好，开心一整天",
      "欢迎 {user}, 上午好，精神满满",
      "欢迎 {user}, 上午好，工作愉快"
    ]
  },
  {
    Enabled: true,
    Key: "noon",
    Random: true,
    Danmu: [
      "欢迎 {user}, 中午好，午餐愉快",
      "欢迎 {user}, 中午好，美食充盈",
      "欢迎 {user}, 中午好，快乐午时",
      "欢迎 {user}, 中午好，记得午间小憩",
      "欢迎 {user}, 中午好，吃好喝好",
      "欢迎 {user}, 中午好，百般顺利"
    ]
  },
  {
    Enabled: true,
    Key: "afternoon",
    Random: true,
    Danmu: [
      "欢迎 {user}, 下午好，动力十足",
      "欢迎 {user}, 下午好，活力满满",
      "欢迎 {user}, 下午好，祝你成功",
      "欢迎 {user}, 下午好，笑容满面",
      "欢迎 {user}, 下午好，开心每时",
      "欢迎 {user}, 下午好，工作顺利",
      "欢迎 {user}, 下午好，美好心情",
      "欢迎 {user}, 下午好，活力四射",
      "欢迎 {user}, 下午好，幸福无限"
    ]
  },
  {
    Enabled: true,
    Key: "night",
    Random: true,
    Danmu: [
      "欢迎 {user}, 晚上好，祝福相伴",
      "欢迎 {user}, 晚上好，幸福安康",
      "欢迎 {user}, 晚上好，心情愉悦",
      "欢迎 {user}, 晚上好，生活美满",
      "欢迎 {user}, 晚上好，夜色温馨",
      "欢迎 {user}, 晚上好，美梦成真",
      "欢迎 {user}, 晚上好，美好的晚上",
      "欢迎 {user}, 祝愿你今晚开心"
    ]
  },
  {
    Enabled: true,
    Key: "midnight",
    Random: true,
    Danmu: [
      "欢迎 {user}, 午夜好，月色清朗",
      "欢迎 {user}, 午夜好，真诚问候",
      "欢迎 {user}, 午夜好，祝福相随",
      "欢迎 {user}, 午夜深沉，你不是独自一人",
      "欢迎 {user}, 午夜好，静夜祝福",
      "欢迎 {user}, 祝福你坚定向前",
      "欢迎 {user}, 午夜沉寂，吉祥相随。",
      "欢迎 {user}, 午夜好，还没休息?",
      "欢迎 {user}, 拥抱午夜的美好"
    ]
  }
];
//凌晨 - Early morning   2:00--5:00
// 早晨 - Morning   5:00--9:00
// 上午 - Late morning / Mid-morning  9:00--11:00
// 中午 - Noon  11:00--14:00
// 下午 - Afternoon 14:00 -- 20:00
// 晚上 - Evening / Night 20:00--00:00
// 午夜 - Midnight 00:00 -- 2:00
const WelcomeDanmuByTimeDescribe = reactive({
  earlymorning: "凌晨2:00-5:00",
  morning: "早晨5:00-9:00",
  latemorning: "上午9:00-11:00",
  noon: "中午11:00-14:00",
  afternoon: "下午14:00-20:00",
  night: "晚上20:00-00:00",
  midnight: "午夜00:00-2:00"
});
import { ref } from "vue";
import type { FormInstance, FormRules, TabsPaneContext } from "element-plus";
import { ReadConfig, WriteConfig } from "../../../wailsjs/go/main/App";
import { ElNotification } from "element-plus";
import { Monitor, Start, Stop, Restart } from "../../../wailsjs/go/main/Program";
// import { func } from "vue-types";

const activeName = ref("first");

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event);
};
const addWelcomeDanmu = () => {
  data.form.WelcomeDanmu.unshift("{user}");
};
const addWelcomeDanmuByTime = index => {
  data.form.WelcomeDanmuByTime[index].Danmu.unshift("{user}");
};
const addWelcomeBlacklist = () => {
  data.form.WelcomeBlacklist.push("");
};
const addWelcomeBlacklistWide = () => {
  data.form.WelcomeBlacklistWide.push("");
};
const deleteWelcomeDanmu = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeDanmu.splice(item, 1);
  }
};
type TierWelcomeKey = "WelcomeNew" | "WelcomeBack" | "WelcomeFan";
const addTierWelcome = (key: TierWelcomeKey) => {
  data.form[key].push("{at}");
};
const deleteTierWelcome = (key: TierWelcomeKey, index: number) => {
  if (index !== -1) {
    data.form[key].splice(index, 1);
  }
};
const deleteWelcomeDanmuByTime = (index1, index2) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (index1 !== -1 && index2 !== -1) {
    data.form.WelcomeDanmuByTime[index1].Danmu.splice(index2, 1);
  }
};
const initWelcomeDanmuByTime = () => {
  if (data.form.WelcomeDanmuByTime == null) {
    data.form.WelcomeDanmuByTime = WelcomeDanmuByTimeTemplate;
    return;
  }
  if (data.form.WelcomeDanmuByTime.length == 0) {
    data.form.WelcomeDanmuByTime = WelcomeDanmuByTimeTemplate;
    return;
  }
  if (data.form.WelcomeDanmuByTime.length < 7) {
    for (let i = 0; i < WelcomeDanmuByTimeTemplate.length; i++) {
      let x = true;
      for (let a = 0; i < data.form.WelcomeDanmuByTime.length; i++) {
        if (
          data.form.WelcomeDanmuByTime[a].Key ===
          WelcomeDanmuByTimeTemplate[i].Key
        ) {
          x = false;
          break;
        }
      }
      if (x) {
        data.form.WelcomeDanmuByTime.push(WelcomeDanmuByTimeTemplate[i]);
      }
    }
  }
  return;
};
const deleteWelcomeBlacklist = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeBlacklist.splice(item, 1);
  }
};
const deleteWelcomeBlacklistWide = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeBlacklistWide.splice(item, 1);
  }
};
const addFocusDanmu = () => {
  data.form.FocusDanmu.push("");
};
const deleteFocusDanmu = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.FocusDanmu.splice(item, 1);
  }
};
const getWelcomeDanmuByTimeDescribe = key => {
  return WelcomeDanmuByTimeDescribe[key];
};

function formatWelcomeString() {
  const arr = Object.entries(data.form.WelcomeString);
  data.tabledata = arr.map(([key, value]) => ({ key, value }));
}

function addWelcomeString() {
  data.dialogdata = {
    uid: "",
    msg: ""
  };
  data.dialogtitle = "新增";
  data.dialogAction = false;
  data.dialogVisible = true;
}

function saveWelcomeString() {
  data.form.WelcomeString[data.dialogdata.uid] = data.dialogdata.msg;
  formatWelcomeString();
  data.dialogVisible = false;
}

function deleteWelcomeString(item: string) {
  // console.log(item)
  if (data.form.WelcomeString.hasOwnProperty(item)) {
    delete data.form.WelcomeString[item];
  }
  formatWelcomeString();
}

function editWelcomeString(uid, msg) {
  data.dialogdata.uid = uid;
  data.dialogdata.msg = msg;
  data.dialogAction = true;
  data.dialogVisible = true;
}

function addDanmulist() {
  data.form.CronDanmuList.push({
    Cron: "*/2 * * * *",
    Random: true,
    Danmu: [""]
  });
}

function deleteDanmulist(row: number) {
  if (row !== -1) {
    data.form.CronDanmuList.splice(row, 1);
  }
}

function addDanmu(item: number) {
  data.form.CronDanmuList[item].Danmu.push("");
}

function deleteDanmu(item: number, row: number) {
  if (item !== -1 && row !== -1) {
    data.form.CronDanmuList[item].Danmu.splice(row, 1);
  }
}

// 关键词回复
function formatKeywordReply() {
  if (data.form.KeywordReplyList) {
    const arr = Object.entries(data.form.KeywordReplyList);
    data.KeywordReplyTemp = arr.map(([key, value]) => ({ key, value }));
  } else {
    data.form.KeywordReplyList = {};
  }
}
function addKeyword() {
  data.keywordDialogData = {
    keyword: "",
    replyword: ""
  };
  data.keywordReplyDialogtitle = "新增";
  data.keywordReplyDialogAction = false;
  data.keywordReplyDialogVisible = true;
}

function saveKeyword() {
  data.form.KeywordReplyList[data.keywordDialogData.keyword] = data.keywordDialogData.replyword;
  formatKeywordReply();
  data.keywordReplyDialogVisible = false;
}

function deleteKeyword(item: string) {
  // console.log(item)
  if (data.form.WelcomeString.hasOwnProperty(item)) {
    delete data.form.KeywordReplyList[item];
  }
  formatKeywordReply();
}
function editKeyword(keyword, replyword) {
  data.keywordReplyDialogtitle = "编辑";
  data.keywordDialogData.keyword = keyword;
  data.keywordDialogData.replyword = replyword;
  data.keywordReplyDialogAction = true;
  data.keywordReplyDialogVisible = true;
}

async function saveConfig() {
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

function restart() {
  Restart();
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
      data.form.WelcomeTier = !!data.form.WelcomeTier;
      data.form.WelcomeNew = data.form.WelcomeNew || [];
      data.form.WelcomeBack = data.form.WelcomeBack || [];
      data.form.WelcomeFan = data.form.WelcomeFan || [];
      data.form.WelcomeFanMedalLevel = data.form.WelcomeFanMedalLevel || 0;
      formatWelcomeString();
      initWelcomeDanmuByTime();
      formatKeywordReply()
      if (data.form.ChatGPT.APIUrl.length == 0) {
        data.form.ChatGPT.APIUrl = "https://api.openai.com/v1";
      }
    }
  });
});
</script>

<template>
  <el-scrollbar class="my-scrollbar">
    <div style="margin: 20px;">
      <el-form :model="data.form" ref="formRef" class="el-form" :rules="rules">
        <el-form-item label="分时间段欢迎">
          <span class="between">
            <el-switch v-model="data.form.InteractWordByTime" />
            <el-tag>Tips: 自定义欢迎语 {user}为用户昵称占位符，列表为随机列表</el-tag>
            <el-button type="primary" v-if="data.form.InteractWordByTime == false"
              @click="addWelcomeDanmu">新增</el-button>
          </span>
        </el-form-item>
        <el-form-item v-if="data.form.InteractWordByTime == false" v-for="(items, index) in data.form.WelcomeDanmu"
          :label="'欢迎弹幕 ' + (index + 1)">
          <div class="between">
            <el-input v-model="data.form.WelcomeDanmu[index]" />
            <el-button class="mt-2" @click="deleteWelcomeDanmu(index)" plain type="danger" :icon="Delete" circle
              style="margin-left: 20px;">
            </el-button>
          </div>
        </el-form-item>
        <el-timeline style="margin-top: 25px;">
          <el-timeline-item :timestamp="'第 ' + (index + 1) + ' 段'" placement="top"
            v-if="data.form.InteractWordByTime == true" v-for="(items, index) in data.form.WelcomeDanmuByTime">
            <el-card>
              <el-form-item>
                <div class="between">
                  <span style="display: flex;">
                    <span style="color: cornflowerblue;font-size: 18px;margin: 0 10px;font-weight: 700;">
                      {{ getWelcomeDanmuByTimeDescribe(data.form.WelcomeDanmuByTime[index].Key) }}
                    </span>
                    <el-switch v-model="data.form.WelcomeDanmuByTime[index].Enabled" />
                  </span>
                  <span>随机弹幕：<el-switch v-model="data.form.WelcomeDanmuByTime[index].Random" /></span>
                  <el-button type="primary" @click="addWelcomeDanmuByTime(index)">新增</el-button>
                </div>
              </el-form-item>
              <el-form-item v-for="(items, index2) in data.form.WelcomeDanmuByTime[index].Danmu"
                :label="'第 ' + (index2 + 1) + ' 条'">
                <div class="between">
                  <el-input v-model="data.form.WelcomeDanmuByTime[index].Danmu[index2]" />
                  <el-button class="mt-2" @click="deleteWelcomeDanmuByTime(index, index2)" plain type="danger"
                    :icon="Delete" circle style="margin-left: 20px;">
                  </el-button>
                </div>
              </el-form-item>
            </el-card>
          </el-timeline-item>
        </el-timeline>
        <el-divider content-position="left">三档欢迎（新人 / 回头客 / 老粉）</el-divider>
        <el-alert type="info" :closable="false" show-icon style="margin-bottom: 16px;">
          <template #title>
            按“老粉 → 新人 → 回头客”顺序匹配；某一档列表为空时回退到上方的随机/分时段欢迎。
            占位符：{at}（@用户，@模式下自动省略）、{anchor}（主播昵称）、{user}（昵称）、{medal}/{medal_level}（粉丝牌）。
            每条建议不超过弹幕最大长度（常见 20 字，部分账号 40 字），超出会被拆成多条。
          </template>
        </el-alert>
        <el-form-item label="启用三档欢迎">
          <el-switch v-model="data.form.WelcomeTier" />
          <el-tag type="info" style="margin-left: 12px;">需同时开启“欢迎弹幕”开关</el-tag>
        </el-form-item>
        <el-form-item label="老粉最低粉丝牌等级">
          <el-input-number v-model="data.form.WelcomeFanMedalLevel" :min="0" :max="60" />
          <span style="margin-left: 12px; color: #909399;">0 表示不限</span>
        </el-form-item>
        <el-card shadow="never" style="margin-bottom: 16px;">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>新人欢迎（首次进入直播间）</span>
              <el-button type="primary" size="small" @click="addTierWelcome('WelcomeNew')">新增</el-button>
            </div>
          </template>
          <el-form-item v-for="(items, index) in data.form.WelcomeNew" :label="'第 ' + (index + 1) + ' 条'">
            <div class="between">
              <el-input v-model="data.form.WelcomeNew[index]" />
              <el-button class="mt-2" @click="deleteTierWelcome('WelcomeNew', index)" plain type="danger"
                :icon="Delete" circle style="margin-left: 20px;" />
            </div>
          </el-form-item>
        </el-card>
        <el-card shadow="never" style="margin-bottom: 16px;">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>回头客欢迎（来过但未佩戴本直播间粉丝牌）</span>
              <el-button type="primary" size="small" @click="addTierWelcome('WelcomeBack')">新增</el-button>
            </div>
          </template>
          <el-form-item v-for="(items, index) in data.form.WelcomeBack" :label="'第 ' + (index + 1) + ' 条'">
            <div class="between">
              <el-input v-model="data.form.WelcomeBack[index]" />
              <el-button class="mt-2" @click="deleteTierWelcome('WelcomeBack', index)" plain type="danger"
                :icon="Delete" circle style="margin-left: 20px;" />
            </div>
          </el-form-item>
        </el-card>
        <el-card shadow="never" style="margin-bottom: 16px;">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>老粉欢迎（佩戴本直播间粉丝牌）</span>
              <el-button type="primary" size="small" @click="addTierWelcome('WelcomeFan')">新增</el-button>
            </div>
          </template>
          <el-form-item v-for="(items, index) in data.form.WelcomeFan" :label="'第 ' + (index + 1) + ' 条'">
            <div class="between">
              <el-input v-model="data.form.WelcomeFan[index]" />
              <el-button class="mt-2" @click="deleteTierWelcome('WelcomeFan', index)" plain type="danger"
                :icon="Delete" circle style="margin-left: 20px;" />
            </div>
          </el-form-item>
        </el-card>
      </el-form>
      <center>
        <el-button type="success" @click="saveConfig" >保存</el-button>
      </center>
    </div>
  </el-scrollbar>
</template>

<style>

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
