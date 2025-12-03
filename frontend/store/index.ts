import { configureStore } from '@reduxjs/toolkit';
import authReducer from './slices/authSlice';
import bettingReducer from './slices/bettingSlice';
import stockReducer from './slices/stockSlice';
import paperTradingReducer from './slices/paperTradingSlice';
import uiReducer from './slices/uiSlice';
import { baseApi } from './api/baseApi';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    betting: bettingReducer,
    stocks: stockReducer,
    paperTrading: paperTradingReducer,
    ui: uiReducer,
    [baseApi.reducerPath]: baseApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(baseApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
