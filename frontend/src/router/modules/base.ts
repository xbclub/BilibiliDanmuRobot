const { VITE_HIDE_HOME } = import.meta.env;
const Layout = () => import("@/layout/index.vue");

export default {
  path: "/configuration",
  name: "Configuration",
  component: Layout,
  redirect: "/base",
  meta: {
    icon: "notification",
    title: "配置页",
    showLink: VITE_HIDE_HOME === "true" ? false : true,
    rank: 0
  },
  children: [
    {
      path: "/base",
      name: "Base",
      component: () => import("@/views/base/index.vue"),
      meta: {
        title: "基础配置",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
    {
      path: "/welcomebarrage",
      name: "Welcomebarrage",
      component: () => import("@/views/welcomeBarrage/index.vue"),
      meta: {
        title: "欢迎弹幕",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
    {
      path: "/keywordreply",
      name: "Keywordreply",
      component: () => import("@/views/keywordReply/index.vue"),
      meta: {
        title: "关键词回复",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
    {
      path: "/assigner",
      name: "Assigner",
      component: () => import("@/views/assigner/index.vue"),
      meta: {
        title: "指定人欢迎",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
    {
      path: "/robots",
      name: "Robots",
      component: () => import("@/views/robots/index.vue"),
      meta: {
        title: "ai机器人",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
    {
      path: "/thank",
      name: "Thank",
      component: () => import("@/views/thank/index.vue"),
      meta: {
        title: "关注答谢语",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
    {
      path: "/timedbarrage",
      name: "Timedbarrage",
      component: () => import("@/views/timedBarrage/index.vue"),
      meta: {
        title: "定时弹幕",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
    {
      path: "/blacklist",
      name: "Blacklist",
      component: () => import("@/views/blackList/index.vue"),
      meta: {
        title: "欢迎黑名单",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    },
  ]
} as RouteConfigsTable;
