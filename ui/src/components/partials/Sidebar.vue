<template>
  <aside class="menu sidebar minimizedsidebar">
    <ul class="menu-list">
      <li
        v-for="(item, index) in menuItems"
        :class="{ current: route.toLowerCase() === item.title.toLowerCase() || (item.children && menuChildernHasActive(item.children) ) }"
        :key="index"
      >
        <router-link
          :to="item.path ? item.path : ''"
          :class="{ 'is-active': item.children && item.active }"
          :exact="true"
          @click.native="toggleAllDropdown(), toggleDropdown(index)"
          v-if="showButton(item.title)"
        >
          <span class="m-icon icon is-small"><i :class="item.icon"></i></span>
          <span class="m-title">{{ item.title }}</span>
          <span
            class="nbComments"
            v-if="item.title === 'Comments' && getHashContext !== '' && getNbComments > 0"
            >{{ getNbComments }}</span
          >
          <span
            class="icon is-small dropdown-icon"
            :class="{ active: item.active }"
            v-if="item.children"
          >
            <i class="ion-ios-arrow-down"></i>
          </span>
        </router-link>
        <ul
          class="dropdown-container"
          v-if="item.children"
          :class="{ active: item.active }"
        >
          <li
            v-for="(child, index) in item.children"
            :class="{
              current: route.toLowerCase() === child.title.toLowerCase() 
            }"
            :key="index"
          >
            <router-link
              class="m-sub-item"
              :to="showButton(child.title) ? child.path : ''"
              :class="{ disabled: !showButton(child.title) }"
            >
              <span class="sub-title">{{ child.title }}</span>
            </router-link>
          </li>
        </ul>
      </li>

      <li ref="toggleSideBar" class="togglebutton">
        <a @click="toggleSideBar()">
          <span class="m-icon icon is-small"><i class="ion-chevron-right"></i></span>
        </a>
      </li>
    </ul>
  </aside>
</template>

<script>
import { mapGetters } from "vuex"
export default {
  data() {
    return {
      init: 0,
      menu: [
        {
          title: "Summary",
          slug: "summary",
          active: false,
          icon: "ion-android-settings",
        },
        {
          title: "Static analysis",
          path: false,
          active: false,
          icon: "ion-stats-bars",
          children: [
            { title: "PE", slug: "pe" },
            {
              title: "Strings",
              slug: "strings",
            },
            {
              title: "Antivirus",
              slug: "antivirus",
            },
          ],
        },
        {
          title: "Comments",
          slug: "comments",
          active: false,
          icon: "ion-chatbubble",
        },
        // {
        //   title: "Dynamic analysis",
        //   path: false,
        //   active: false,
        //   icon: "ion-ios-analytics",
        //   children: [
        //     { title: "API Calls", path: "/" },
        //     { title: "Network", path: "/" },
        //     { title: "Dropped files", path: "/" },
        //     { title: "Memory dumps", path: "/" },
        //   ],
        // }, 
      ],
    }
  },
  computed: {
    ...mapGetters(["getHashContext", "getNbComments", "getLoggedIn", "isPE"]),
   
    menuItems: function() {
      const hash = this.getHashContext
      return this.menu.map(({ slug, children, ...item }) => ({
        ...item,
        ...(slug ? { path: `/${slug}/${hash}` } : {}),
        ...(children
          ? {
              children: children.map(({ slug, ...item }) => ({
                ...item,
                ...(slug ? { path: `/${slug}/${hash}` } : {}),
              })),
            }
          : {}),
      }))
    },
    route: function() {
      const routeName = this.$route.path.replace(/\/\//g, "/").split('/').filter(x=> x)[0] || "";
      return routeName
    },
  },
  mounted() {
    this.menu.map((el) => {
      el.active = this.menuChildernHasActive(el.children)
    })
    // this.menu[0].active = false
    this.menu.map((el) => {
      el.dropdownHeight = el.active ? el.children.length * 36 + "px" : 0
    });
  },
  methods: {
    menuChildernHasActive: function(children){
      return children && children.reduce((x,v)=> x || this.route.toLowerCase() === v.title.toLowerCase(), false);
    },
    toggleAllDropdown(){
      this.menu.forEach((x,i)=> this.toggleDropdown(i, false));
    },
    toggleSideBar(){
      document.querySelector('#app').classList.toggle('minimized');
      this.$refs.toggleSideBar.querySelector('i').classList.toggle('ion-chevron-right');
      this.$refs.toggleSideBar.querySelector('i').classList.toggle('ion-chevron-left');
    },
    toggleDropdown(index, toggle=true) {
      if (this.menu[index] && this.menu[index].children) {
        let found = !toggle
        this.menu.forEach((el, i) => {
          if (el.active && i === index) {
            el.active = false
            el.dropdownHeight = 0
            found = true
          }
        })
        if (found || !toggle) return
        // this.menu.map((el) => {el.active = false; el.dropdownHeight = 0})
        this.menu[index].active = true
        this.menu[index].dropdownHeight =
          this.menu[index].children.length * 36 + "px"
      }
    },
    showButton(name) {
      if (name !== "Comments" && name !== "PE") return true
      else if (name === "Comments" && this.getHashContext !== "") return true
      else if (name === "PE" && this.isPE) return true
      else return false
    },
  }
}
</script>

<style scoped lang="scss">
@import "../../assets/scss/variables";
aside.sidebar {
  z-index: 101;
  background-color: #fff;
  position: fixed;
  top: $header-height;
  left: 0;
  height: calc(100% - #{$header-height});
  width: 200px;
  border-right: 1px solid rgba(25,25,25,0.1);

  *{
    transition: padding 100ms linear;
    transition: margin 100ms linear;
  }
  
  .menu-list {
    .is-active {
      background-color: transparent;
      color: #4a4a4a;
    }

    .dropdown-icon {
      float: right;
      line-height: 2rem;

      i {
        transition: all 0.2s;
      }

      &.active {
        i {
          transform: rotate(180deg);
        }
      }
    }

    li{
      padding: 3px 0;
    }

    ul.dropdown-container {
      transition:all 5s;

      position: absolute;
      height: 0;
      overflow: hidden;
      transition: all 0.2s;
      padding: 0;
      margin: 0;
      background: white;

      &.active {
        height: auto;
      }
      .disabled {
        cursor: not-allowed;
        opacity: 0.5;
        text-decoration: none;
      }
    }
    ul.dropdown-container.active{
      position: unset;
    }
    li:not(.current) > a:hover {
      background-color: #fff;
      color: $primary-color;
    }

    li.current > a:hover {
      background-color: $primary-color;
    }
  }
  li .m-icon{
    margin-right: 0.5rem;
  }
}
.sidebar:not(.minimizedsidebar){
  .nbComments {
    position: absolute;
    border: 2px solid #00e0bf;
    background-color: #ffffff;
    color: black;
    display: inline-block;
    border-radius: 50%;
    min-width: 1.3rem;
    height: 1.3rem;
    line-height: 1rem;
    font-size: 0.8rem;
    margin-top: 27px;
    margin-left: 4px;
  }
  
}
.minimized {
  
  .nbComments {
    position: absolute;
    border: 2px solid #00e0bfa6;
    background-color: #ffffff;
    color: black;
    display: inline-block;
    border-radius: 50%;
    min-width: 1.4rem;
    height: 1.4rem;
    line-height: 1rem;
    font-size: 0.9rem;
    margin-top: 25px;
    margin-left: 5px;
  }
  .minimizedsidebar {
    max-width: $sidebar-width;
    span {
      text-align: center;
    }
    a .m-title{
      display: none;
    }
    a .m-icon {
      display: block;
    }

    a:not(.m-sub-item)  {
      margin: 13px;
      border-radius: 7px;
      padding: 12px 0 6px 0px;
    }

    .menu-list{
      li {
        text-align: center;
        line-height: 2rem;
      }
      .m-icon{
        display: contents;
        text-align: center;
        font-size: 2rem;
      }
      .dropdown-icon{
        display: none;
      }

      li .m-sub-item{
        text-align: left;
        padding: 0px 8px;
      }
    }
    .current {
      background-color: unset !important;
      .m-icon{
        background: $primary-color;
      }
      > a{
        background: lighten($primary-color, 3%) !important;
        color:white !important;
      }
    }

    .is-active > .current a:first-child{
        background: lighten($primary-color, 3%) !important;
        color:white;
    }

  }
}
.togglebutton{
    text-align: center;
    padding: 20px 0;
    font-size: 1.3rem;
    bottom: 0;
    position: absolute;
    width: 100%;
}

</style>
