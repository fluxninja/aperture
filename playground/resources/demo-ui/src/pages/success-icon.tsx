import { SvgIcon, SvgIconProps } from '@mui/material'
import React, { FC, PropsWithRef } from 'react'

export const SuccessIcon: FC<PropsWithRef<SvgIconProps>> = React.forwardRef(
  (props, ref) => (
    <SvgIcon width="100" height="100" viewBox="0 0 49 44" {...props} ref={ref}>
      <path
        d="M39.8854 23.1555C39.2407 22.2571 39.3436 21.0021 40.1391 20.234C40.7015 19.6854 41.353 19.1985 42.0525 18.7801C42.752 18.3618 43.1841 17.628 43.1841 16.8119V14.1784C43.1841 6.97058 36.1889 6.29163 32.4375 5.92129C27.2734 5.42066 25.9224 4.81715 25.3052 2.8969C25.1955 2.54714 25.0103 2.2111 24.7428 1.95735C22.4454 -0.237225 19.373 0.839488 18.502 3.13693C17.8436 4.87201 16.3966 5.44123 11.452 5.92129C7.69375 6.28477 0.705405 6.96372 0.705405 14.1784V16.8187C0.705405 17.628 1.14432 18.3686 1.83698 18.787C2.52964 19.2053 3.18116 19.6922 3.74352 20.2409C4.53905 21.009 4.64192 22.264 3.99726 23.1624C3.48291 23.8756 2.79711 24.438 2.00843 24.7878C1.21976 25.1375 0.698547 25.9605 0.698547 26.8452V29.8284C0.698547 37.0362 7.69375 37.7152 11.4451 38.0855C16.6092 38.5861 17.9602 39.1896 18.5775 41.1099C18.6872 41.4597 18.8724 41.7957 19.1398 42.0494C21.4373 44.244 24.5097 43.1673 25.3806 40.8699C26.039 39.1348 27.486 38.5656 32.4307 38.0855C36.1889 37.722 43.1772 37.0431 43.1772 29.8284V26.8452C43.1772 25.9605 42.6766 25.1512 41.8674 24.7878C41.0581 24.4243 40.3929 23.8756 39.8785 23.1624L39.8854 23.1555Z"
        fill="#EFEEED"
      />
      <path
        d="M48.7322 18.3138C46.0096 18.1835 42.8549 19.5688 41.1198 21.6742C41.8262 23.423 43.6641 24.4037 45.509 23.9922C44.9809 22.6686 43.8013 21.7976 42.464 21.6193C44.8232 21.3382 47.2715 20.0763 48.7322 18.3138Z"
        fill="#F8773D"
      />
      <path
        d="M13.1116 19.7609C12.7824 19.7609 12.4669 19.8294 12.1858 19.9597C12.3846 20.1449 12.4875 20.4467 12.4052 20.7553C12.3229 21.0502 12.0623 21.2696 11.7606 21.3108C11.4862 21.3451 11.2394 21.2491 11.0748 21.0639C10.8827 21.4891 10.8142 21.9829 10.9239 22.4972C11.1296 23.4231 11.932 24.1363 12.8784 24.2323C14.2226 24.3695 15.3542 23.3133 15.3542 22.0034C15.3542 20.7621 14.3529 19.7609 13.1116 19.7609ZM14.1197 23.6974C13.8248 23.6974 13.5917 23.4642 13.5917 23.1693C13.5917 22.8744 13.8248 22.6412 14.1197 22.6412C14.4146 22.6412 14.6478 22.8813 14.6478 23.1693C14.6478 23.4573 14.4078 23.6974 14.1197 23.6974Z"
        fill="#56AE89"
      />
      <path
        d="M28.7685 21.0638C28.5696 21.4959 28.501 21.9965 28.6245 22.5177C28.8165 23.3476 29.4817 24.0128 30.3184 24.1911C31.8889 24.5272 33.2605 23.2104 33.0274 21.6536C32.8765 20.6661 32.0604 19.8843 31.066 19.7745C30.6408 19.7265 30.2361 19.7951 29.8795 19.9597C30.0784 20.1449 30.1881 20.4329 30.1058 20.7415C30.0304 21.0158 29.7972 21.2422 29.516 21.297C29.2349 21.3519 28.9468 21.2559 28.7685 21.0638ZM31.2854 23.1624C31.2854 22.8744 31.5254 22.6343 31.8135 22.6343C32.1015 22.6343 32.3484 22.8744 32.3484 23.1624C32.3484 23.4504 32.1084 23.6905 31.8135 23.6905C31.5186 23.6905 31.2854 23.4573 31.2854 23.1624Z"
        fill="#56AE89"
      />
      <path
        d="M31.1003 16.3112C26.9786 16.3112 28.3639 19.2465 21.9448 19.9803C16.3623 19.3425 16.6847 17.0313 14.1403 16.4415C10.7456 15.6597 7.34401 18.1355 7.1177 21.6125C6.89138 25.0895 9.52487 27.6887 12.7962 27.6887C16.911 27.6887 15.5394 24.7466 21.9448 24.0197C27.5341 24.6575 27.2049 26.9686 29.7492 27.5515C33.144 28.3334 36.5455 25.8576 36.7719 22.3806C36.9913 19.0682 34.3647 16.3044 31.0934 16.3044M13.1116 25.3364C11.2394 25.3364 9.72375 23.7865 9.77176 21.9005C9.81977 20.1586 11.2737 18.7047 13.0156 18.6567C14.9016 18.6018 16.4515 20.1174 16.4515 21.9965C16.4515 23.8756 14.9564 25.3364 13.1116 25.3364ZM30.7779 25.3364C28.8988 25.3364 27.3832 23.7865 27.4381 21.9005C27.4861 20.1586 28.94 18.7047 30.6888 18.6567C32.5747 18.6087 34.1178 20.1243 34.1178 21.9965C34.1178 23.8688 32.6227 25.3364 30.7779 25.3364Z"
        fill="#56AE89"
      />
      <path
        d="M39.5082 17.6074C38.7332 17.6074 38.1023 16.9764 38.1023 16.2015V14.1783C38.1023 11.9221 36.7169 11.4488 31.9438 10.9825C28.5079 10.6465 24.3793 10.2487 21.9447 7.25173C19.5101 10.2487 15.3816 10.6465 11.9457 10.9825C7.17938 11.4488 5.7872 11.9221 5.7872 14.1783V16.2015C5.7872 16.9764 5.15626 17.6074 4.3813 17.6074C3.60634 17.6074 2.9754 16.9764 2.9754 16.2015V14.1783C2.9754 9.03482 7.7966 8.56847 11.6714 8.19128C16.28 7.74551 19.3935 7.20372 20.628 3.94615C20.8337 3.39751 21.3549 3.04089 21.9447 3.04089C22.5277 3.04089 23.0489 3.40437 23.2615 3.94615C24.4959 7.20372 27.6095 7.74551 32.2181 8.19128C36.0929 8.56847 40.9141 9.03482 40.9141 14.1783V16.2015C40.9141 16.9764 40.2831 17.6074 39.5082 17.6074Z"
        fill="#F8773D"
      />
      <path
        d="M21.9447 40.9659C21.3618 40.9659 20.8406 40.6024 20.628 40.0606C19.3935 36.8031 16.28 36.2681 11.6714 35.8155C7.7966 35.4383 2.9754 34.972 2.9754 29.8284V27.8053C2.9754 27.0304 3.60634 26.3994 4.3813 26.3994C5.15626 26.3994 5.7872 27.0304 5.7872 27.8053V29.8284C5.7872 32.0847 7.17252 32.5579 11.9457 33.0243C15.3816 33.3603 19.5101 33.7581 21.9447 36.7551C24.3793 33.7581 28.5079 33.3603 31.9438 33.0243C36.7101 32.5579 38.1023 32.0847 38.1023 29.8284V27.8053C38.1023 27.0304 38.7332 26.3994 39.5082 26.3994C40.2831 26.3994 40.9141 27.0304 40.9141 27.8053V29.8284C40.9141 34.972 36.0929 35.4383 32.2181 35.8155C27.6095 36.2613 24.4959 36.8031 23.2615 40.0606C23.0557 40.6093 22.5345 40.9659 21.9447 40.9659Z"
        fill="#F8773D"
      />
    </SvgIcon>
  )
)

SuccessIcon.displayName = 'SuccessIcon'
