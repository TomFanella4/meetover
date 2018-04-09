import { Toast } from 'native-base';

export const StyledToast = config => (
  Toast.show({
    ...config,
    textStyle: { fontFamily: 'pt-sans' },
  })
);
