import cv2
import rest

import numpy as np


class ChromaKeyServiceImpl(rest.ChromaKeyingService):
    def replace(self, src_image_str, bg_image_str) -> bytes:
        bg = cv2.imdecode(np.frombuffer(bg_image_str, np.uint8), cv2.IMREAD_COLOR)
        img = cv2.imdecode(np.frombuffer(src_image_str, np.uint8), cv2.IMREAD_COLOR)

        RED, GREEN, BLUE = (2, 1, 0)

        reds = img[:, :, RED]
        greens = img[:, :, GREEN]
        blues = img[:, :, BLUE]

        # z = np.zeros(shape=img.shape, dtype=in

        mask = (greens < 70) | (reds > greens) | (blues > greens)
        mask = mask.astype("uint8") * 255

        # print(mask)

        mask_inv = cv2.bitwise_not(mask)

        # cv2.imshow("Mask", mask)
        # cv2.imshow("Mask inv", mask_inv)

        # converting mask 2d to 3d
        result = cv2.bitwise_and(img, img, mask=mask)

        bg = cv2.resize(bg, (1280, 720))
        bg = cv2.bitwise_and(bg, bg, mask=mask_inv)

        res = cv2.add(result, bg)

        is_success, im_buf_arr = cv2.imencode(".jpg", res)
        return im_buf_arr.tobytes()

        # cv2.imshow("Result", res)
        # # cv2.imshow("Bg", bg)
        # cv2.waitKey(0)
        # cv2.destroyAllWindows()
