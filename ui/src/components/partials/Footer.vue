<template>
  <footer class="main-footer" :class="{ fullwidth: fullwidth }">
    <div class="container">
      <div class="columns">
        <!-- <div class="column is-4">
                <h3>Join the developers</h3>
                <p>
                    Install in seconds on Linux or OS XXx.
                </p>
                <a href="" class="button secondary-color-btn sm-btn bold-btn">Install</a>
                <a href="" class="button transparent-white-btn sm-btn bold-btn">Learn more</a>
            </div>  -->
        <div class="column is-3">
          <ul class="footer-links">
            <li>
              <header>getting started</header>
            </li>
            <li>
              <a href="">faq</a>
            </li>
            <li>
              <a href="">install</a>
            </li>
            <li>
              <a href="">quick start</a>
            </li>
            <li>
              <a href="">tutorial</a>
            </li>
          </ul>
        </div>
        <div class="column is-3">
          <ul class="footer-links">
            <li>
              <header>documentation</header>
            </li>
            <li>
              <a href="">api</a>
            </li>
          </ul>
        </div>
        <div class="column is-3">
          <ul class="footer-links">
            <li>
              <header>community</header>
            </li>
            <li>
              <a href="">events</a>
            </li>
            <li>
              <a href="">github</a>
            </li>
            <li>
              <a href="">twitter</a>
            </li>
          </ul>
        </div>
        <div class="column is-3">
          <ul class="footer-links">
            <li>
              <header>saferwall</header>
            </li>
            <li>
              <a href="">blog</a>
            </li>
          </ul>
        </div>
      </div>
      <!-- <div class="columns">
            <div class="column is-3">
                <div class="base">
                    <i class="icon ion-chatbubbles"></i>
                    <span>Template based on
                        <br>
                        <a href="">Rethinkdb.com</a>
                    </span>
                </div>
            </div>
            <div class="column is-8 is-offset-1">
                <div class="github-footer">
                    <p>
                        <i class="icon ion-social-github"></i>
                        <a href="">Star this project</a>
                        on Github.
                    </p>
                </div>
            </div>
        </div> -->
    </div>
  </footer>
</template>
<script>
export default {
  props: ["fullwidth"],
}
</script>
<style lang="scss" scoped>
footer.main-footer {
  background-color: #12282f;
  padding: 50px 0 120px;
  color: #fff;
  float: right;
  width: 100%;
  margin-top: 150px;

  @media screen and (max-width: 468px) {
    text-align: center;
  }

  h3 {
    line-height: normal;
    font-size: 15px;
  }

  p {
    font-size: 14px;
  }

  .footer-links {
    li {
      line-height: 24px;

      header {
        padding-bottom: 5px;
        margin-bottom: 10px;
        font-weight: 600;
        border-bottom: solid 1px #1a3741;
        display: inline-block;
      }

      a {
        color: rgba(white, 0.8);
        font-size: 13px;
      }
    }
  }

  .github-footer {
    p {
      font-size: 22px;
      color: rgba(white, 0.7);
      line-height: 25px;
      float: left;
    }

    .icon {
      font-size: 30px;
      display: inline-block;
      margin-right: 10px;
      color: rgba(white, 0.7);
    }

    a {
      font-size: 22px;
      border-bottom: dotted 1px rgba(white, 0.25);
      color: rgba(white, 0.7);
      display: inline-block;
      padding-bottom: 5px;
    }
  }

  .base {
    border-top: solid 1px #1a3741;
    padding-top: 30px;
    line-height: 20px;
    font-size: 13px;

    &,
    * {
      color: rgba(white, 0.7);
    }

    .icon {
      float: left;
      font-size: 24px;
      color: rgba(white, 0.5);
    }

    span {
      padding-left: 20px;
      float: left;
    }
  }
}
</style>
