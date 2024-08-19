import TaskPage from "@/interfaces/TaskPage";
import DownloadPage from "@/interfaces/DownloadPage";
import SettingPage from "@/interfaces/SettingPage";
import React from "react";

const routes: Array<[string, () => React.JSX.Element]> = [];

function EmptyPage() {
  return <div></div>;
}

const Router = {
  add(route: string, constructor: (...args: any) => React.JSX.Element) {
    routes.push([route, constructor]);
  },
  search: function (searchRoute: string): [(...args: any) => React.JSX.Element, { [key: string]: any }] {
    for (const [route, constructor] of routes) {
      if (searchRoute === route) {
        return [constructor, {}];
      }
      const searchComponents = searchRoute.split("/");
      const routeComponents = route.split("/");
      if (searchComponents.length !== routeComponents.length) {
        continue;
      }
        
      const params: { [key: string]: any } = {};
      let match = true;

      for (let i = 0; i < searchComponents.length; i++) {
        if (searchComponents[i] !== routeComponents[i] && !routeComponents[i].startsWith(":")) {
          match = false;
          break;
        }
        if (routeComponents[i].startsWith(":")) {
          params[routeComponents[i].substring(1)] = searchComponents[i];
        }
      }
      if (match) {
        return [constructor, params];
      }
    }
    return [EmptyPage, {}];
  },
  route: function (searchRoute: string, data: any): React.JSX.Element {
    const [constructor, _params] = this.search(searchRoute);
    return React.createElement(constructor, data);
  },
};

Router.add('/task', TaskPage);
Router.add('/download', DownloadPage);
Router.add('/setting', SettingPage);

export default Router;
