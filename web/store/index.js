export const state = () => {
    return {
        connectInfo: null,
        isConnect: false,
        errorLogs: [],
    }
  }

  export const actions = {
    saveConnectInfo({commit},params){
        commit('SAVE_CONNECT_INFO',params)
    },
    saveConnectState({commit},params){
        commit("SAVE_CONNECT_STATE",params)
    },
    saveErrorLogs({commit},params){
        commit('SAVE_ERROR_LOGS',params)
    },
  }

  export const mutations = {
    SAVE_CONNECT_INFO(state,params){
        state.connectInfo = params
    },
    SAVE_CONNECT_STATE(state,params){
        state.isConnect = params
    },
    SAVE_ERROR_LOGS(state,params){
        state.errorLogs = params
    }
  }

  export const getters = {
    "getConnectInfo": state => state.connectInfo,
    "getConnectState": state => state.isConnect,
    "getErrorLogs": state => state.errorLogs,
  }
