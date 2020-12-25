<template>
  <aside class="menu sidebar modern">
    <ul class="menu-list">
      <li
        v-for="(item, index) in menuItems"
        :class="{ current: route.toLowerCase() === item.title.toLowerCase() }"
        :key="index"
      >
        <router-link
          :to="item.path ? item.path : ''"
          :class="{ 'is-active': item.children && item.active }"
          :exact="true"
          @click.native="toggleDropdown(index)"
          v-if="showButton(item.title)"
        >
          <span class="m-icon icon is-small"><i :class="item.icon"></i></span>
          <span class="m-title">{{ item.title }}</span>
          <span
            class="nbComments"
            v-if="item.title === 'Comments' && getHashContext !== ''"
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
    </ul>
  </aside>
</template>

<script>
import { mapGetters } from "vuex"
export default {
  data() {
    return {
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
      el.active = Boolean(el.children)
    })
    // this.menu[0].active = false
    this.menu.map((el) => {
      el.dropdownHeight = el.active ? el.children.length * 36 + "px" : 0
    })
  },
  methods: {
    toggleDropdown(index) {
      if (this.menu[index].children) {
        let found = false
        this.menu.forEach((el, i) => {
          if (el.active && i === index) {
            el.active = false
            el.dropdownHeight = 0
            found = true
          }
        })
        if (found) return
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
  },
}
</script>

<style scoped lang="scss">
@import "../../assets/scss/variables";
aside.sidebar {
  background-color: #fff;
  position: fixed;
  top: $header-height;
  left: 0;
  height: calc(100% - #{$header-height});
  width: 200px;
  box-shadow: 0 0 30px rgba(black, 0.05);
  padding-top: 15px;

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

    .current {
      background-color: $primary-color;
    }

    ul.dropdown-container {
      height: 0;
      overflow: hidden;
      transition: all 0.2s;
      padding-left: 0;
      margin: 0 0 0 0.75em;

      &.active {
        height: auto;
        margin: 0 0 0.75em 0.75em;
      }
      .disabled {
        cursor: not-allowed;
        opacity: 0.5;
        text-decoration: none;
      }
    }

    li:not(.current) > a:hover {
      background-color: #fff;
      color: $primary-color;
    }

    li.current > a:hover {
      background-color: $primary-color;
    }
  }
}
.nbComments {
  position: absolute;
  border: 2px solid #00e0bfa6;
  background-color: #ffffff;
  color: black;
  display: inline-block;
  border-radius: 50%;
  min-width: 1.9em;
  height: 1.9em;
  font-size: 0.7rem;
  font-weight: 600;
  line-height: 1rem;
  margin-top: 25px;
  margin-left: 5px;
}
.modern {
  background-color: $primary-color;
  transition: max-width 300ms;
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
    a{
      background: lighten($primary-color, 3%) !important;
      color:white;
    }
  }

  .is-active > .current{
      background: lighten($primary-color, 3%) !important;
      color:white;
  }

}

</style>
