import cv2
import numpy as np
import logging

log = logging.getLogger("bluefin-tuna")


def preprocess(img):
    if len(img.shape) == 3 and img.shape[2] == 3:
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    elif len(img.shape) == 3 and img.shape[2] == 4:
        img = img[:, :, 0:3]
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    elif len(img.shape) == 2:
        gray = img
    else:
        raise ValueError("wrong input image!")
    blur = cv2.GaussianBlur(gray, (3, 3), 0)
    return blur


def key_points_color(gray_img, drift=50):
    h, w = gray_img.shape[0], gray_img.shape[1]
    if w <= 10 or h <= 10:
        raise ValueError("image size ({}, {}) too small!".format(w, h))
    if h < 100:
        h_bound = 3
    else:
        h_bound = 5
    if w < 100:
        w_bound = 3
    else:
        w_bound = 5

    key_points = [np.mean(gray_img[0:h_bound, 0:w_bound]),
                  np.mean(gray_img[h - h_bound - 1:-1, 0:w_bound]),
                  np.mean(gray_img[0:h_bound, w - w_bound - 1:-1]),
                  np.mean(gray_img[h - h_bound - 1:-1, w - w_bound - 1:-1])]
    log.debug(key_points)
    if max(key_points) - min(key_points) < drift:
        return True, sum(key_points) / len(key_points)
    else:
        return False, 0


def main_color(gray_img):
    bins = 256
    h, w = gray_img.shape[0], gray_img.shape[1]
    pixel_num = h * w
    color_map = {}
    for i in range(bins):
        color_map[i] = len(np.where(gray_img == i)[0]) / pixel_num
    max_color = -1
    max_ratio = 0
    for k, v in color_map.items():
        if v > max_ratio:
            max_color = k
            max_ratio = v
    return max_color, max_ratio


def generate_mask(img, th=0.55, drift=20):
    gray = preprocess(img)
    mask = np.ones(gray.shape) * 255
    is_same, key_color = key_points_color(gray)
    max_color, max_ratio = main_color(gray)
    log.debug("max_ratio: {}".format(max_ratio))
    if not is_same:
        log.debug("Warning: key points color not same!")
        raise ValueError("key points color not same!")
    if key_color != max_color:
        log.debug("Warning: key color != max color!")
    if max_ratio > th:
        mask[((key_color - int(drift / 2)) <= gray) & (gray <= (key_color + int(drift / 2)))] = 0
        mask[((max_color - int(drift / 2)) <= gray) & (gray <= (max_color + int(drift / 2)))] = 0
    else:
        mask[gray == key_color] = 0
    return mask


def blend(img, mask):
    if img.shape[0:2] != mask.shape:
        raise ValueError("image shape != mask shape: image shape: {}, mask shape:{}".format(img.shape[0:2], mask.shape))
    rgba = np.zeros((mask.shape[0], mask.shape[1], 4), dtype=int)
    rgba[:, :, 0:3] = img[:, :, 0:3]
    rgba[:, :, 3] = mask
    return rgba