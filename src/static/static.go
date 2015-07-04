package static

import (
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path"
	"path/filepath"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

// static_1block_bin reads file data from disk. It returns an error on failure.
func static_1block_bin() (*asset, error) {
	path := "/home/z/IdeaProjects/static/1block.bin"
	name := "static/1block.bin"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_as3cam_css reads file data from disk. It returns an error on failure.
func static_css_as3cam_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/AS3Cam.css"
	name := "static/css/AS3Cam.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_jcrop_gif reads file data from disk. It returns an error on failure.
func static_css_jcrop_gif() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/Jcrop.gif"
	name := "static/css/Jcrop.gif"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_bootstrap_responsive_css reads file data from disk. It returns an error on failure.
func static_css_bootstrap_responsive_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/bootstrap-responsive.css"
	name := "static/css/bootstrap-responsive.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_bootstrap_css reads file data from disk. It returns an error on failure.
func static_css_bootstrap_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/bootstrap.css"
	name := "static/css/bootstrap.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_bootstrap_min_css reads file data from disk. It returns an error on failure.
func static_css_bootstrap_min_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/bootstrap.min.css"
	name := "static/css/bootstrap.min.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_cf_css reads file data from disk. It returns an error on failure.
func static_css_cf_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/cf.css"
	name := "static/css/cf.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_cropper_css reads file data from disk. It returns an error on failure.
func static_css_cropper_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/cropper.css"
	name := "static/css/cropper.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_font_awesome_css reads file data from disk. It returns an error on failure.
func static_css_font_awesome_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/font-awesome.css"
	name := "static/css/font-awesome.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_jquery_ui_timepicker_addon_css reads file data from disk. It returns an error on failure.
func static_css_jquery_ui_timepicker_addon_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/jquery-ui-timepicker-addon.css"
	name := "static/css/jquery-ui-timepicker-addon.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_jquery_ui_css reads file data from disk. It returns an error on failure.
func static_css_jquery_ui_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/jquery-ui.css"
	name := "static/css/jquery-ui.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_jquery_jcrop_css reads file data from disk. It returns an error on failure.
func static_css_jquery_jcrop_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/jquery.Jcrop.css"
	name := "static/css/jquery.Jcrop.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_jquery_jcrop_min_css reads file data from disk. It returns an error on failure.
func static_css_jquery_jcrop_min_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/jquery.Jcrop.min.css"
	name := "static/css/jquery.Jcrop.min.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_jquery_css reads file data from disk. It returns an error on failure.
func static_css_jquery_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/jquery.css"
	name := "static/css/jquery.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_jquery_qtip_min_css reads file data from disk. It returns an error on failure.
func static_css_jquery_qtip_min_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/jquery.qtip.min.css"
	name := "static/css/jquery.qtip.min.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_plugins_metismenu_metismenu_min_css reads file data from disk. It returns an error on failure.
func static_css_plugins_metismenu_metismenu_min_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/plugins/metisMenu/metisMenu.min.css"
	name := "static/css/plugins/metisMenu/metisMenu.min.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_progress_css reads file data from disk. It returns an error on failure.
func static_css_progress_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/progress.css"
	name := "static/css/progress.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_sb_admin_2_css reads file data from disk. It returns an error on failure.
func static_css_sb_admin_2_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/sb-admin-2.css"
	name := "static/css/sb-admin-2.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_sb_admin_css reads file data from disk. It returns an error on failure.
func static_css_sb_admin_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/sb-admin.css"
	name := "static/css/sb-admin.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_social_buttons_css reads file data from disk. It returns an error on failure.
func static_css_social_buttons_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/social-buttons.css"
	name := "static/css/social-buttons.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_tooltipster_shadow_css reads file data from disk. It returns an error on failure.
func static_css_tooltipster_shadow_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/tooltipster-shadow.css"
	name := "static/css/tooltipster-shadow.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_tooltipster_css reads file data from disk. It returns an error on failure.
func static_css_tooltipster_css() (*asset, error) {
	path := "/home/z/IdeaProjects/static/css/tooltipster.css"
	name := "static/css/tooltipster.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_otf reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_otf() (*asset, error) {
	path := "/home/z/IdeaProjects/static/fonts/FontAwesome.otf"
	name := "static/fonts/FontAwesome.otf"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_eot reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_eot() (*asset, error) {
	path := "/home/z/IdeaProjects/static/fonts/fontawesome-webfont.eot"
	name := "static/fonts/fontawesome-webfont.eot"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_svg reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_svg() (*asset, error) {
	path := "/home/z/IdeaProjects/static/fonts/fontawesome-webfont.svg"
	name := "static/fonts/fontawesome-webfont.svg"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_ttf reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_ttf() (*asset, error) {
	path := "/home/z/IdeaProjects/static/fonts/fontawesome-webfont.ttf"
	name := "static/fonts/fontawesome-webfont.ttf"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_woff reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_woff() (*asset, error) {
	path := "/home/z/IdeaProjects/static/fonts/fontawesome-webfont.woff"
	name := "static/fonts/fontawesome-webfont.woff"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_luxisr_ttf reads file data from disk. It returns an error on failure.
func static_fonts_luxisr_ttf() (*asset, error) {
	path := "/home/z/IdeaProjects/static/fonts/luxisr.ttf"
	name := "static/fonts/luxisr.ttf"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_alert_png reads file data from disk. It returns an error on failure.
func static_img_alert_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/alert.png"
	name := "static/img/alert.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_blank_png reads file data from disk. It returns an error on failure.
func static_img_blank_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/blank.png"
	name := "static/img/blank.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_cf_blurb_img_png reads file data from disk. It returns an error on failure.
func static_img_cf_blurb_img_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/cf_blurb_img.png"
	name := "static/img/cf_blurb_img.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_face_jpg reads file data from disk. It returns an error on failure.
func static_img_face_jpg() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/face.jpg"
	name := "static/img/face.jpg"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_glyphicons_halflings_png reads file data from disk. It returns an error on failure.
func static_img_glyphicons_halflings_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/glyphicons-halflings.png"
	name := "static/img/glyphicons-halflings.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_k_bg_png reads file data from disk. It returns an error on failure.
func static_img_k_bg_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/k_bg.png"
	name := "static/img/k_bg.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_k_bg_pass_png reads file data from disk. It returns an error on failure.
func static_img_k_bg_pass_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/k_bg_pass.png"
	name := "static/img/k_bg_pass.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_logo_small_png reads file data from disk. It returns an error on failure.
func static_img_logo_small_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/logo-small.png"
	name := "static/img/logo-small.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_logo_png reads file data from disk. It returns an error on failure.
func static_img_logo_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/logo.png"
	name := "static/img/logo.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_noavatar_png reads file data from disk. It returns an error on failure.
func static_img_noavatar_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/noavatar.png"
	name := "static/img/noavatar.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_photo_png reads file data from disk. It returns an error on failure.
func static_img_photo_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/photo.png"
	name := "static/img/photo.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_profile_jpg reads file data from disk. It returns an error on failure.
func static_img_profile_jpg() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/profile.jpg"
	name := "static/img/profile.jpg"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_race_gif reads file data from disk. It returns an error on failure.
func static_img_race_gif() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/race.gif"
	name := "static/img/race.gif"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_us_ru_png reads file data from disk. It returns an error on failure.
func static_img_us_ru_png() (*asset, error) {
	path := "/home/z/IdeaProjects/static/img/us-ru.png"
	name := "static/img/us-ru.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_serpent_js reads file data from disk. It returns an error on failure.
func static_js_serpent_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/Serpent.js"
	name := "static/js/Serpent.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_aes_js reads file data from disk. It returns an error on failure.
func static_js_aes_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/aes.js"
	name := "static/js/aes.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_asn1hex_1_1_js reads file data from disk. It returns an error on failure.
func static_js_asn1hex_1_1_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/asn1hex-1.1.js"
	name := "static/js/asn1hex-1.1.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_asn1hex_1_1_min_js reads file data from disk. It returns an error on failure.
func static_js_asn1hex_1_1_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/asn1hex-1.1.min.js"
	name := "static/js/asn1hex-1.1.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_asn1hex_1_js reads file data from disk. It returns an error on failure.
func static_js_asn1hex_1_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/asn1hex-1.js"
	name := "static/js/asn1hex-1.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_base64_js reads file data from disk. It returns an error on failure.
func static_js_base64_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/base64.js"
	name := "static/js/base64.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_alert_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_alert_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-alert.js"
	name := "static/js/bootstrap-alert.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_button_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_button_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-button.js"
	name := "static/js/bootstrap-button.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_carousel_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_carousel_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-carousel.js"
	name := "static/js/bootstrap-carousel.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_collapse_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_collapse_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-collapse.js"
	name := "static/js/bootstrap-collapse.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_dropdown_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_dropdown_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-dropdown.js"
	name := "static/js/bootstrap-dropdown.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_modal_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_modal_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-modal.js"
	name := "static/js/bootstrap-modal.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_popover_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_popover_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-popover.js"
	name := "static/js/bootstrap-popover.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_scrollspy_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_scrollspy_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-scrollspy.js"
	name := "static/js/bootstrap-scrollspy.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_tab_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_tab_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-tab.js"
	name := "static/js/bootstrap-tab.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_tooltip_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_tooltip_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-tooltip.js"
	name := "static/js/bootstrap-tooltip.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_transition_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_transition_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-transition.js"
	name := "static/js/bootstrap-transition.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_typeahead_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_typeahead_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap-typeahead.js"
	name := "static/js/bootstrap-typeahead.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_bootstrap_min_js reads file data from disk. It returns an error on failure.
func static_js_bootstrap_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/bootstrap.min.js"
	name := "static/js/bootstrap.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_cropper_js reads file data from disk. It returns an error on failure.
func static_js_cropper_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/cropper.js"
	name := "static/js/cropper.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_enc_base64_min_js reads file data from disk. It returns an error on failure.
func static_js_enc_base64_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/enc-base64-min.js"
	name := "static/js/enc-base64-min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_html5shiv_js reads file data from disk. It returns an error on failure.
func static_js_html5shiv_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/html5shiv.js"
	name := "static/js/html5shiv.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_index_js reads file data from disk. It returns an error on failure.
func static_js_index_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/index.js"
	name := "static/js/index.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_infobubble_js reads file data from disk. It returns an error on failure.
func static_js_infobubble_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/infobubble.js"
	name := "static/js/infobubble.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_1_11_0_js reads file data from disk. It returns an error on failure.
func static_js_jquery_1_11_0_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery-1.11.0.js"
	name := "static/js/jquery-1.11.0.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_1_9_1_min_js reads file data from disk. It returns an error on failure.
func static_js_jquery_1_9_1_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery-1.9.1.min.js"
	name := "static/js/jquery-1.9.1.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_ui_slideraccess_js reads file data from disk. It returns an error on failure.
func static_js_jquery_ui_slideraccess_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery-ui-sliderAccess.js"
	name := "static/js/jquery-ui-sliderAccess.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_ui_timepicker_addon_js reads file data from disk. It returns an error on failure.
func static_js_jquery_ui_timepicker_addon_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery-ui-timepicker-addon.js"
	name := "static/js/jquery-ui-timepicker-addon.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_ui_min_js reads file data from disk. It returns an error on failure.
func static_js_jquery_ui_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery-ui.min.js"
	name := "static/js/jquery-ui.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_jcrop_js reads file data from disk. It returns an error on failure.
func static_js_jquery_jcrop_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery.Jcrop.js"
	name := "static/js/jquery.Jcrop.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_js reads file data from disk. It returns an error on failure.
func static_js_jquery_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery.js"
	name := "static/js/jquery.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_min_js reads file data from disk. It returns an error on failure.
func static_js_jquery_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery.min.js"
	name := "static/js/jquery.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_qtip_min_js reads file data from disk. It returns an error on failure.
func static_js_jquery_qtip_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery.qtip.min.js"
	name := "static/js/jquery.qtip.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_tooltipster_min_js reads file data from disk. It returns an error on failure.
func static_js_jquery_tooltipster_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery.tooltipster.min.js"
	name := "static/js/jquery.tooltipster.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_webcam_as3_js reads file data from disk. It returns an error on failure.
func static_js_jquery_webcam_as3_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery.webcam.as3.js"
	name := "static/js/jquery.webcam.as3.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_002_js reads file data from disk. It returns an error on failure.
func static_js_jquery_002_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jquery_002.js"
	name := "static/js/jquery_002.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_js_js reads file data from disk. It returns an error on failure.
func static_js_js_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/js.js"
	name := "static/js/js.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jsbn_js reads file data from disk. It returns an error on failure.
func static_js_jsbn_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jsbn.js"
	name := "static/js/jsbn.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jsbn2_js reads file data from disk. It returns an error on failure.
func static_js_jsbn2_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/jsbn2.js"
	name := "static/js/jsbn2.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_markerclusterer_js reads file data from disk. It returns an error on failure.
func static_js_markerclusterer_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/markerclusterer.js"
	name := "static/js/markerclusterer.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_mcrypt_js reads file data from disk. It returns an error on failure.
func static_js_mcrypt_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/mcrypt.js"
	name := "static/js/mcrypt.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_md5_js reads file data from disk. It returns an error on failure.
func static_js_md5_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/md5.js"
	name := "static/js/md5.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_plugins_metismenu_metismenu_min_js reads file data from disk. It returns an error on failure.
func static_js_plugins_metismenu_metismenu_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/plugins/metisMenu/metisMenu.min.js"
	name := "static/js/plugins/metisMenu/metisMenu.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_prng4_js reads file data from disk. It returns an error on failure.
func static_js_prng4_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/prng4.js"
	name := "static/js/prng4.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rijndael_js reads file data from disk. It returns an error on failure.
func static_js_rijndael_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rijndael.js"
	name := "static/js/rijndael.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_ripemd160_js reads file data from disk. It returns an error on failure.
func static_js_ripemd160_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/ripemd160.js"
	name := "static/js/ripemd160.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rng_js reads file data from disk. It returns an error on failure.
func static_js_rng_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rng.js"
	name := "static/js/rng.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsa_js reads file data from disk. It returns an error on failure.
func static_js_rsa_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsa.js"
	name := "static/js/rsa.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsa2_js reads file data from disk. It returns an error on failure.
func static_js_rsa2_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsa2.js"
	name := "static/js/rsa2.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsapem_1_1_js reads file data from disk. It returns an error on failure.
func static_js_rsapem_1_1_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsapem-1.1.js"
	name := "static/js/rsapem-1.1.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsapem_1_1_min_js reads file data from disk. It returns an error on failure.
func static_js_rsapem_1_1_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsapem-1.1.min.js"
	name := "static/js/rsapem-1.1.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsapem_1_js reads file data from disk. It returns an error on failure.
func static_js_rsapem_1_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsapem-1.js"
	name := "static/js/rsapem-1.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsasign_1_2_js reads file data from disk. It returns an error on failure.
func static_js_rsasign_1_2_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsasign-1.2.js"
	name := "static/js/rsasign-1.2.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsasign_1_2_min_js reads file data from disk. It returns an error on failure.
func static_js_rsasign_1_2_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsasign-1.2.min.js"
	name := "static/js/rsasign-1.2.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_rsasign_1_js reads file data from disk. It returns an error on failure.
func static_js_rsasign_1_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/rsasign-1.js"
	name := "static/js/rsasign-1.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_sas3cam_js reads file data from disk. It returns an error on failure.
func static_js_sas3cam_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/sAS3Cam.js"
	name := "static/js/sAS3Cam.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_sb_admin_2_js reads file data from disk. It returns an error on failure.
func static_js_sb_admin_2_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/sb-admin-2.js"
	name := "static/js/sb-admin-2.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_sb_admin_js reads file data from disk. It returns an error on failure.
func static_js_sb_admin_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/sb-admin.js"
	name := "static/js/sb-admin.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_sha1_js reads file data from disk. It returns an error on failure.
func static_js_sha1_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/sha1.js"
	name := "static/js/sha1.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_sha256_js reads file data from disk. It returns an error on failure.
func static_js_sha256_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/sha256.js"
	name := "static/js/sha256.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_sha512_js reads file data from disk. It returns an error on failure.
func static_js_sha512_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/sha512.js"
	name := "static/js/sha512.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_spin_js reads file data from disk. It returns an error on failure.
func static_js_spin_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/spin.js"
	name := "static/js/spin.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_spots_js reads file data from disk. It returns an error on failure.
func static_js_spots_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/spots.js"
	name := "static/js/spots.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_stacktable_js reads file data from disk. It returns an error on failure.
func static_js_stacktable_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/stacktable.js"
	name := "static/js/stacktable.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_unixtime_js reads file data from disk. It returns an error on failure.
func static_js_unixtime_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/unixtime.js"
	name := "static/js/unixtime.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_uploader_js reads file data from disk. It returns an error on failure.
func static_js_uploader_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/uploader.js"
	name := "static/js/uploader.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_x509_1_1_js reads file data from disk. It returns an error on failure.
func static_js_x509_1_1_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/x509-1.1.js"
	name := "static/js/x509-1.1.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_x509_1_1_min_js reads file data from disk. It returns an error on failure.
func static_js_x509_1_1_min_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/x509-1.1.min.js"
	name := "static/js/x509-1.1.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_youtube_webcam_js reads file data from disk. It returns an error on failure.
func static_js_youtube_webcam_js() (*asset, error) {
	path := "/home/z/IdeaProjects/static/js/youtube_webcam.js"
	name := "static/js/youtube_webcam.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_lang_1_ini reads file data from disk. It returns an error on failure.
func static_lang_1_ini() (*asset, error) {
	path := "/home/z/IdeaProjects/static/lang/1.ini"
	name := "static/lang/1.ini"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_lang_42_ini reads file data from disk. It returns an error on failure.
func static_lang_42_ini() (*asset, error) {
	path := "/home/z/IdeaProjects/static/lang/42.ini"
	name := "static/lang/42.ini"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_lang_en_us_all_json reads file data from disk. It returns an error on failure.
func static_lang_en_us_all_json() (*asset, error) {
	path := "/home/z/IdeaProjects/static/lang/en-us.all.json"
	name := "static/lang/en-us.all.json"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_lang_locale_en_us_ini reads file data from disk. It returns an error on failure.
func static_lang_locale_en_us_ini() (*asset, error) {
	path := "/home/z/IdeaProjects/static/lang/locale_en-US.ini"
	name := "static/lang/locale_en-US.ini"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_lang_locale_ru_ru_ini reads file data from disk. It returns an error on failure.
func static_lang_locale_ru_ru_ini() (*asset, error) {
	path := "/home/z/IdeaProjects/static/lang/locale_ru-RU.ini"
	name := "static/lang/locale_ru-RU.ini"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_nodes_inc reads file data from disk. It returns an error on failure.
func static_nodes_inc() (*asset, error) {
	path := "/home/z/IdeaProjects/static/nodes.inc"
	name := "static/nodes.inc"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_abuse_html reads file data from disk. It returns an error on failure.
func static_templates_abuse_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/abuse.html"
	name := "static/templates/abuse.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_add_cf_project_data_html reads file data from disk. It returns an error on failure.
func static_templates_add_cf_project_data_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/add_cf_project_data.html"
	name := "static/templates/add_cf_project_data.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_after_install_html reads file data from disk. It returns an error on failure.
func static_templates_after_install_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/after_install.html"
	name := "static/templates/after_install.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_alert_success_html reads file data from disk. It returns an error on failure.
func static_templates_alert_success_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/alert_success.html"
	name := "static/templates/alert_success.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_arbitration_html reads file data from disk. It returns an error on failure.
func static_templates_arbitration_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/arbitration.html"
	name := "static/templates/arbitration.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_arbitration_arbitrator_html reads file data from disk. It returns an error on failure.
func static_templates_arbitration_arbitrator_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/arbitration_arbitrator.html"
	name := "static/templates/arbitration_arbitrator.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_arbitration_buyer_html reads file data from disk. It returns an error on failure.
func static_templates_arbitration_buyer_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/arbitration_buyer.html"
	name := "static/templates/arbitration_buyer.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_arbitration_seller_html reads file data from disk. It returns an error on failure.
func static_templates_arbitration_seller_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/arbitration_seller.html"
	name := "static/templates/arbitration_seller.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_assignments_html reads file data from disk. It returns an error on failure.
func static_templates_assignments_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/assignments.html"
	name := "static/templates/assignments.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_assignments_new_miner_html reads file data from disk. It returns an error on failure.
func static_templates_assignments_new_miner_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/assignments_new_miner.html"
	name := "static/templates/assignments_new_miner.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_assignments_promised_amount_html reads file data from disk. It returns an error on failure.
func static_templates_assignments_promised_amount_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/assignments_promised_amount.html"
	name := "static/templates/assignments_promised_amount.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_block_explorer_html reads file data from disk. It returns an error on failure.
func static_templates_block_explorer_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/block_explorer.html"
	name := "static/templates/block_explorer.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_bug_reporting_html reads file data from disk. It returns an error on failure.
func static_templates_bug_reporting_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/bug_reporting.html"
	name := "static/templates/bug_reporting.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_cash_request_in_html reads file data from disk. It returns an error on failure.
func static_templates_cash_request_in_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/cash_request_in.html"
	name := "static/templates/cash_request_in.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_cash_request_out_html reads file data from disk. It returns an error on failure.
func static_templates_cash_request_out_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/cash_request_out.html"
	name := "static/templates/cash_request_out.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_cf_catalog_html reads file data from disk. It returns an error on failure.
func static_templates_cf_catalog_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/cf_catalog.html"
	name := "static/templates/cf_catalog.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_cf_page_preview_html reads file data from disk. It returns an error on failure.
func static_templates_cf_page_preview_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/cf_page_preview.html"
	name := "static/templates/cf_page_preview.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_cf_project_change_category_html reads file data from disk. It returns an error on failure.
func static_templates_cf_project_change_category_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/cf_project_change_category.html"
	name := "static/templates/cf_project_change_category.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_cf_start_html reads file data from disk. It returns an error on failure.
func static_templates_cf_start_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/cf_start.html"
	name := "static/templates/cf_start.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_arbitrator_conditions_html reads file data from disk. It returns an error on failure.
func static_templates_change_arbitrator_conditions_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_arbitrator_conditions.html"
	name := "static/templates/change_arbitrator_conditions.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_avatar_tpl reads file data from disk. It returns an error on failure.
func static_templates_change_avatar_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_avatar.tpl"
	name := "static/templates/change_avatar.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_commission_html reads file data from disk. It returns an error on failure.
func static_templates_change_commission_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_commission.html"
	name := "static/templates/change_commission.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_country_race_tpl reads file data from disk. It returns an error on failure.
func static_templates_change_country_race_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_country_race.tpl"
	name := "static/templates/change_country_race.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_creditor_html reads file data from disk. It returns an error on failure.
func static_templates_change_creditor_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_creditor.html"
	name := "static/templates/change_creditor.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_geolocation_html reads file data from disk. It returns an error on failure.
func static_templates_change_geolocation_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_geolocation.html"
	name := "static/templates/change_geolocation.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_host_html reads file data from disk. It returns an error on failure.
func static_templates_change_host_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_host.html"
	name := "static/templates/change_host.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_key_close_html reads file data from disk. It returns an error on failure.
func static_templates_change_key_close_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_key_close.html"
	name := "static/templates/change_key_close.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_key_request_html reads file data from disk. It returns an error on failure.
func static_templates_change_key_request_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_key_request.html"
	name := "static/templates/change_key_request.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_money_back_time_html reads file data from disk. It returns an error on failure.
func static_templates_change_money_back_time_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_money_back_time.html"
	name := "static/templates/change_money_back_time.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_node_key_html reads file data from disk. It returns an error on failure.
func static_templates_change_node_key_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_node_key.html"
	name := "static/templates/change_node_key.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_primary_key_html reads file data from disk. It returns an error on failure.
func static_templates_change_primary_key_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_primary_key.html"
	name := "static/templates/change_primary_key.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_change_promised_amount_html reads file data from disk. It returns an error on failure.
func static_templates_change_promised_amount_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/change_promised_amount.html"
	name := "static/templates/change_promised_amount.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_credits_html reads file data from disk. It returns an error on failure.
func static_templates_credits_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/credits.html"
	name := "static/templates/credits.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_currency_exchange_html reads file data from disk. It returns an error on failure.
func static_templates_currency_exchange_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/currency_exchange.html"
	name := "static/templates/currency_exchange.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_currency_exchange_delete_html reads file data from disk. It returns an error on failure.
func static_templates_currency_exchange_delete_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/currency_exchange_delete.html"
	name := "static/templates/currency_exchange_delete.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_db_info_html reads file data from disk. It returns an error on failure.
func static_templates_db_info_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/db_info.html"
	name := "static/templates/db_info.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_del_cf_funding_html reads file data from disk. It returns an error on failure.
func static_templates_del_cf_funding_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/del_cf_funding.html"
	name := "static/templates/del_cf_funding.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_del_cf_project_html reads file data from disk. It returns an error on failure.
func static_templates_del_cf_project_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/del_cf_project.html"
	name := "static/templates/del_cf_project.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_del_credit_html reads file data from disk. It returns an error on failure.
func static_templates_del_credit_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/del_credit.html"
	name := "static/templates/del_credit.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_del_promised_amount_html reads file data from disk. It returns an error on failure.
func static_templates_del_promised_amount_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/del_promised_amount.html"
	name := "static/templates/del_promised_amount.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_for_repaid_fix_html reads file data from disk. It returns an error on failure.
func static_templates_for_repaid_fix_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/for_repaid_fix.html"
	name := "static/templates/for_repaid_fix.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_holidays_list_html reads file data from disk. It returns an error on failure.
func static_templates_holidays_list_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/holidays_list.html"
	name := "static/templates/holidays_list.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_home_html reads file data from disk. It returns an error on failure.
func static_templates_home_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/home.html"
	name := "static/templates/home.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_index_html reads file data from disk. It returns an error on failure.
func static_templates_index_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/index.html"
	name := "static/templates/index.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_index_cf_html reads file data from disk. It returns an error on failure.
func static_templates_index_cf_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/index_cf.html"
	name := "static/templates/index_cf.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_information_html reads file data from disk. It returns an error on failure.
func static_templates_information_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/information.html"
	name := "static/templates/information.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_install_step_0_html reads file data from disk. It returns an error on failure.
func static_templates_install_step_0_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/install_step_0.html"
	name := "static/templates/install_step_0.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_install_step_1_html reads file data from disk. It returns an error on failure.
func static_templates_install_step_1_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/install_step_1.html"
	name := "static/templates/install_step_1.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_install_step_2_html reads file data from disk. It returns an error on failure.
func static_templates_install_step_2_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/install_step_2.html"
	name := "static/templates/install_step_2.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_interface_html reads file data from disk. It returns an error on failure.
func static_templates_interface_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/interface.html"
	name := "static/templates/interface.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_list_cf_projects_tpl reads file data from disk. It returns an error on failure.
func static_templates_list_cf_projects_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/list_cf_projects.tpl"
	name := "static/templates/list_cf_projects.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_login_html reads file data from disk. It returns an error on failure.
func static_templates_login_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/login.html"
	name := "static/templates/login.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_menu_html reads file data from disk. It returns an error on failure.
func static_templates_menu_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/menu.html"
	name := "static/templates/menu.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_mining_menu_html reads file data from disk. It returns an error on failure.
func static_templates_mining_menu_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/mining_menu.html"
	name := "static/templates/mining_menu.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_mining_promised_amount_html reads file data from disk. It returns an error on failure.
func static_templates_mining_promised_amount_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/mining_promised_amount.html"
	name := "static/templates/mining_promised_amount.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_modal_html reads file data from disk. It returns an error on failure.
func static_templates_modal_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/modal.html"
	name := "static/templates/modal.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_money_back_html reads file data from disk. It returns an error on failure.
func static_templates_money_back_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/money_back.html"
	name := "static/templates/money_back.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_money_back_request_html reads file data from disk. It returns an error on failure.
func static_templates_money_back_request_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/money_back_request.html"
	name := "static/templates/money_back_request.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_my_cf_projects_html reads file data from disk. It returns an error on failure.
func static_templates_my_cf_projects_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/my_cf_projects.html"
	name := "static/templates/my_cf_projects.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_new_cf_project_html reads file data from disk. It returns an error on failure.
func static_templates_new_cf_project_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/new_cf_project.html"
	name := "static/templates/new_cf_project.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_new_credit_html reads file data from disk. It returns an error on failure.
func static_templates_new_credit_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/new_credit.html"
	name := "static/templates/new_credit.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_new_holidays_html reads file data from disk. It returns an error on failure.
func static_templates_new_holidays_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/new_holidays.html"
	name := "static/templates/new_holidays.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_new_promised_amount_html reads file data from disk. It returns an error on failure.
func static_templates_new_promised_amount_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/new_promised_amount.html"
	name := "static/templates/new_promised_amount.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_new_user_html reads file data from disk. It returns an error on failure.
func static_templates_new_user_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/new_user.html"
	name := "static/templates/new_user.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_node_config_html reads file data from disk. It returns an error on failure.
func static_templates_node_config_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/node_config.html"
	name := "static/templates/node_config.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_notifications_html reads file data from disk. It returns an error on failure.
func static_templates_notifications_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/notifications.html"
	name := "static/templates/notifications.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_pct_tpl reads file data from disk. It returns an error on failure.
func static_templates_pct_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/pct.tpl"
	name := "static/templates/pct.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_points_html reads file data from disk. It returns an error on failure.
func static_templates_points_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/points.html"
	name := "static/templates/points.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_pool_admin_html reads file data from disk. It returns an error on failure.
func static_templates_pool_admin_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/pool_admin.html"
	name := "static/templates/pool_admin.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_pool_tech_works_tpl reads file data from disk. It returns an error on failure.
func static_templates_pool_tech_works_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/pool_tech_works.tpl"
	name := "static/templates/pool_tech_works.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_progress_tpl reads file data from disk. It returns an error on failure.
func static_templates_progress_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/progress.tpl"
	name := "static/templates/progress.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_progress_bar_html reads file data from disk. It returns an error on failure.
func static_templates_progress_bar_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/progress_bar.html"
	name := "static/templates/progress_bar.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_promised_amount_actualization_tpl reads file data from disk. It returns an error on failure.
func static_templates_promised_amount_actualization_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/promised_amount_actualization.tpl"
	name := "static/templates/promised_amount_actualization.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_promised_amount_list_html reads file data from disk. It returns an error on failure.
func static_templates_promised_amount_list_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/promised_amount_list.html"
	name := "static/templates/promised_amount_list.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_reduction_tpl reads file data from disk. It returns an error on failure.
func static_templates_reduction_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/reduction.tpl"
	name := "static/templates/reduction.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_repayment_credit_html reads file data from disk. It returns an error on failure.
func static_templates_repayment_credit_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/repayment_credit.html"
	name := "static/templates/repayment_credit.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_restoring_access_html reads file data from disk. It returns an error on failure.
func static_templates_restoring_access_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/restoring_access.html"
	name := "static/templates/restoring_access.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_rewrite_primary_key_html reads file data from disk. It returns an error on failure.
func static_templates_rewrite_primary_key_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/rewrite_primary_key.html"
	name := "static/templates/rewrite_primary_key.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_sign_up_in_the_pool_html reads file data from disk. It returns an error on failure.
func static_templates_sign_up_in_the_pool_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/sign_up_in_the_pool.html"
	name := "static/templates/sign_up_in_the_pool.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_signatures_html reads file data from disk. It returns an error on failure.
func static_templates_signatures_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/signatures.html"
	name := "static/templates/signatures.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_statistic_html reads file data from disk. It returns an error on failure.
func static_templates_statistic_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/statistic.html"
	name := "static/templates/statistic.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_statistic_voting_html reads file data from disk. It returns an error on failure.
func static_templates_statistic_voting_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/statistic_voting.html"
	name := "static/templates/statistic_voting.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_updating_blockchain_html reads file data from disk. It returns an error on failure.
func static_templates_updating_blockchain_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/updating_blockchain.html"
	name := "static/templates/updating_blockchain.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade.html"
	name := "static/templates/upgrade.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_0_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_0_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_0.html"
	name := "static/templates/upgrade_0.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_1_and_2_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_1_and_2_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_1_and_2.html"
	name := "static/templates/upgrade_1_and_2.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_3_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_3_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_3.html"
	name := "static/templates/upgrade_3.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_4_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_4_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_4.html"
	name := "static/templates/upgrade_4.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_5_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_5_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_5.html"
	name := "static/templates/upgrade_5.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_6_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_6_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_6.html"
	name := "static/templates/upgrade_6.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_7_html reads file data from disk. It returns an error on failure.
func static_templates_upgrade_7_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_7.html"
	name := "static/templates/upgrade_7.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_upgrade_resend_tpl reads file data from disk. It returns an error on failure.
func static_templates_upgrade_resend_tpl() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/upgrade_resend.tpl"
	name := "static/templates/upgrade_resend.tpl"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_vote_for_me_html reads file data from disk. It returns an error on failure.
func static_templates_vote_for_me_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/vote_for_me.html"
	name := "static/templates/vote_for_me.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_voting_html reads file data from disk. It returns an error on failure.
func static_templates_voting_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/voting.html"
	name := "static/templates/voting.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_templates_wallets_list_html reads file data from disk. It returns an error on failure.
func static_templates_wallets_list_html() (*asset, error) {
	path := "/home/z/IdeaProjects/static/templates/wallets_list.html"
	name := "static/templates/wallets_list.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"static/1block.bin": static_1block_bin,
	"static/css/AS3Cam.css": static_css_as3cam_css,
	"static/css/Jcrop.gif": static_css_jcrop_gif,
	"static/css/bootstrap-responsive.css": static_css_bootstrap_responsive_css,
	"static/css/bootstrap.css": static_css_bootstrap_css,
	"static/css/bootstrap.min.css": static_css_bootstrap_min_css,
	"static/css/cf.css": static_css_cf_css,
	"static/css/cropper.css": static_css_cropper_css,
	"static/css/font-awesome.css": static_css_font_awesome_css,
	"static/css/jquery-ui-timepicker-addon.css": static_css_jquery_ui_timepicker_addon_css,
	"static/css/jquery-ui.css": static_css_jquery_ui_css,
	"static/css/jquery.Jcrop.css": static_css_jquery_jcrop_css,
	"static/css/jquery.Jcrop.min.css": static_css_jquery_jcrop_min_css,
	"static/css/jquery.css": static_css_jquery_css,
	"static/css/jquery.qtip.min.css": static_css_jquery_qtip_min_css,
	"static/css/plugins/metisMenu/metisMenu.min.css": static_css_plugins_metismenu_metismenu_min_css,
	"static/css/progress.css": static_css_progress_css,
	"static/css/sb-admin-2.css": static_css_sb_admin_2_css,
	"static/css/sb-admin.css": static_css_sb_admin_css,
	"static/css/social-buttons.css": static_css_social_buttons_css,
	"static/css/tooltipster-shadow.css": static_css_tooltipster_shadow_css,
	"static/css/tooltipster.css": static_css_tooltipster_css,
	"static/fonts/FontAwesome.otf": static_fonts_fontawesome_otf,
	"static/fonts/fontawesome-webfont.eot": static_fonts_fontawesome_webfont_eot,
	"static/fonts/fontawesome-webfont.svg": static_fonts_fontawesome_webfont_svg,
	"static/fonts/fontawesome-webfont.ttf": static_fonts_fontawesome_webfont_ttf,
	"static/fonts/fontawesome-webfont.woff": static_fonts_fontawesome_webfont_woff,
	"static/fonts/luxisr.ttf": static_fonts_luxisr_ttf,
	"static/img/alert.png": static_img_alert_png,
	"static/img/blank.png": static_img_blank_png,
	"static/img/cf_blurb_img.png": static_img_cf_blurb_img_png,
	"static/img/face.jpg": static_img_face_jpg,
	"static/img/glyphicons-halflings.png": static_img_glyphicons_halflings_png,
	"static/img/k_bg.png": static_img_k_bg_png,
	"static/img/k_bg_pass.png": static_img_k_bg_pass_png,
	"static/img/logo-small.png": static_img_logo_small_png,
	"static/img/logo.png": static_img_logo_png,
	"static/img/noavatar.png": static_img_noavatar_png,
	"static/img/photo.png": static_img_photo_png,
	"static/img/profile.jpg": static_img_profile_jpg,
	"static/img/race.gif": static_img_race_gif,
	"static/img/us-ru.png": static_img_us_ru_png,
	"static/js/Serpent.js": static_js_serpent_js,
	"static/js/aes.js": static_js_aes_js,
	"static/js/asn1hex-1.1.js": static_js_asn1hex_1_1_js,
	"static/js/asn1hex-1.1.min.js": static_js_asn1hex_1_1_min_js,
	"static/js/asn1hex-1.js": static_js_asn1hex_1_js,
	"static/js/base64.js": static_js_base64_js,
	"static/js/bootstrap-alert.js": static_js_bootstrap_alert_js,
	"static/js/bootstrap-button.js": static_js_bootstrap_button_js,
	"static/js/bootstrap-carousel.js": static_js_bootstrap_carousel_js,
	"static/js/bootstrap-collapse.js": static_js_bootstrap_collapse_js,
	"static/js/bootstrap-dropdown.js": static_js_bootstrap_dropdown_js,
	"static/js/bootstrap-modal.js": static_js_bootstrap_modal_js,
	"static/js/bootstrap-popover.js": static_js_bootstrap_popover_js,
	"static/js/bootstrap-scrollspy.js": static_js_bootstrap_scrollspy_js,
	"static/js/bootstrap-tab.js": static_js_bootstrap_tab_js,
	"static/js/bootstrap-tooltip.js": static_js_bootstrap_tooltip_js,
	"static/js/bootstrap-transition.js": static_js_bootstrap_transition_js,
	"static/js/bootstrap-typeahead.js": static_js_bootstrap_typeahead_js,
	"static/js/bootstrap.min.js": static_js_bootstrap_min_js,
	"static/js/cropper.js": static_js_cropper_js,
	"static/js/enc-base64-min.js": static_js_enc_base64_min_js,
	"static/js/html5shiv.js": static_js_html5shiv_js,
	"static/js/index.js": static_js_index_js,
	"static/js/infobubble.js": static_js_infobubble_js,
	"static/js/jquery-1.11.0.js": static_js_jquery_1_11_0_js,
	"static/js/jquery-1.9.1.min.js": static_js_jquery_1_9_1_min_js,
	"static/js/jquery-ui-sliderAccess.js": static_js_jquery_ui_slideraccess_js,
	"static/js/jquery-ui-timepicker-addon.js": static_js_jquery_ui_timepicker_addon_js,
	"static/js/jquery-ui.min.js": static_js_jquery_ui_min_js,
	"static/js/jquery.Jcrop.js": static_js_jquery_jcrop_js,
	"static/js/jquery.js": static_js_jquery_js,
	"static/js/jquery.min.js": static_js_jquery_min_js,
	"static/js/jquery.qtip.min.js": static_js_jquery_qtip_min_js,
	"static/js/jquery.tooltipster.min.js": static_js_jquery_tooltipster_min_js,
	"static/js/jquery.webcam.as3.js": static_js_jquery_webcam_as3_js,
	"static/js/jquery_002.js": static_js_jquery_002_js,
	"static/js/js.js": static_js_js_js,
	"static/js/jsbn.js": static_js_jsbn_js,
	"static/js/jsbn2.js": static_js_jsbn2_js,
	"static/js/markerclusterer.js": static_js_markerclusterer_js,
	"static/js/mcrypt.js": static_js_mcrypt_js,
	"static/js/md5.js": static_js_md5_js,
	"static/js/plugins/metisMenu/metisMenu.min.js": static_js_plugins_metismenu_metismenu_min_js,
	"static/js/prng4.js": static_js_prng4_js,
	"static/js/rijndael.js": static_js_rijndael_js,
	"static/js/ripemd160.js": static_js_ripemd160_js,
	"static/js/rng.js": static_js_rng_js,
	"static/js/rsa.js": static_js_rsa_js,
	"static/js/rsa2.js": static_js_rsa2_js,
	"static/js/rsapem-1.1.js": static_js_rsapem_1_1_js,
	"static/js/rsapem-1.1.min.js": static_js_rsapem_1_1_min_js,
	"static/js/rsapem-1.js": static_js_rsapem_1_js,
	"static/js/rsasign-1.2.js": static_js_rsasign_1_2_js,
	"static/js/rsasign-1.2.min.js": static_js_rsasign_1_2_min_js,
	"static/js/rsasign-1.js": static_js_rsasign_1_js,
	"static/js/sAS3Cam.js": static_js_sas3cam_js,
	"static/js/sb-admin-2.js": static_js_sb_admin_2_js,
	"static/js/sb-admin.js": static_js_sb_admin_js,
	"static/js/sha1.js": static_js_sha1_js,
	"static/js/sha256.js": static_js_sha256_js,
	"static/js/sha512.js": static_js_sha512_js,
	"static/js/spin.js": static_js_spin_js,
	"static/js/spots.js": static_js_spots_js,
	"static/js/stacktable.js": static_js_stacktable_js,
	"static/js/unixtime.js": static_js_unixtime_js,
	"static/js/uploader.js": static_js_uploader_js,
	"static/js/x509-1.1.js": static_js_x509_1_1_js,
	"static/js/x509-1.1.min.js": static_js_x509_1_1_min_js,
	"static/js/youtube_webcam.js": static_js_youtube_webcam_js,
	"static/lang/1.ini": static_lang_1_ini,
	"static/lang/42.ini": static_lang_42_ini,
	"static/lang/en-us.all.json": static_lang_en_us_all_json,
	"static/lang/locale_en-US.ini": static_lang_locale_en_us_ini,
	"static/lang/locale_ru-RU.ini": static_lang_locale_ru_ru_ini,
	"static/nodes.inc": static_nodes_inc,
	"static/templates/abuse.html": static_templates_abuse_html,
	"static/templates/add_cf_project_data.html": static_templates_add_cf_project_data_html,
	"static/templates/after_install.html": static_templates_after_install_html,
	"static/templates/alert_success.html": static_templates_alert_success_html,
	"static/templates/arbitration.html": static_templates_arbitration_html,
	"static/templates/arbitration_arbitrator.html": static_templates_arbitration_arbitrator_html,
	"static/templates/arbitration_buyer.html": static_templates_arbitration_buyer_html,
	"static/templates/arbitration_seller.html": static_templates_arbitration_seller_html,
	"static/templates/assignments.html": static_templates_assignments_html,
	"static/templates/assignments_new_miner.html": static_templates_assignments_new_miner_html,
	"static/templates/assignments_promised_amount.html": static_templates_assignments_promised_amount_html,
	"static/templates/block_explorer.html": static_templates_block_explorer_html,
	"static/templates/bug_reporting.html": static_templates_bug_reporting_html,
	"static/templates/cash_request_in.html": static_templates_cash_request_in_html,
	"static/templates/cash_request_out.html": static_templates_cash_request_out_html,
	"static/templates/cf_catalog.html": static_templates_cf_catalog_html,
	"static/templates/cf_page_preview.html": static_templates_cf_page_preview_html,
	"static/templates/cf_project_change_category.html": static_templates_cf_project_change_category_html,
	"static/templates/cf_start.html": static_templates_cf_start_html,
	"static/templates/change_arbitrator_conditions.html": static_templates_change_arbitrator_conditions_html,
	"static/templates/change_avatar.tpl": static_templates_change_avatar_tpl,
	"static/templates/change_commission.html": static_templates_change_commission_html,
	"static/templates/change_country_race.tpl": static_templates_change_country_race_tpl,
	"static/templates/change_creditor.html": static_templates_change_creditor_html,
	"static/templates/change_geolocation.html": static_templates_change_geolocation_html,
	"static/templates/change_host.html": static_templates_change_host_html,
	"static/templates/change_key_close.html": static_templates_change_key_close_html,
	"static/templates/change_key_request.html": static_templates_change_key_request_html,
	"static/templates/change_money_back_time.html": static_templates_change_money_back_time_html,
	"static/templates/change_node_key.html": static_templates_change_node_key_html,
	"static/templates/change_primary_key.html": static_templates_change_primary_key_html,
	"static/templates/change_promised_amount.html": static_templates_change_promised_amount_html,
	"static/templates/credits.html": static_templates_credits_html,
	"static/templates/currency_exchange.html": static_templates_currency_exchange_html,
	"static/templates/currency_exchange_delete.html": static_templates_currency_exchange_delete_html,
	"static/templates/db_info.html": static_templates_db_info_html,
	"static/templates/del_cf_funding.html": static_templates_del_cf_funding_html,
	"static/templates/del_cf_project.html": static_templates_del_cf_project_html,
	"static/templates/del_credit.html": static_templates_del_credit_html,
	"static/templates/del_promised_amount.html": static_templates_del_promised_amount_html,
	"static/templates/for_repaid_fix.html": static_templates_for_repaid_fix_html,
	"static/templates/holidays_list.html": static_templates_holidays_list_html,
	"static/templates/home.html": static_templates_home_html,
	"static/templates/index.html": static_templates_index_html,
	"static/templates/index_cf.html": static_templates_index_cf_html,
	"static/templates/information.html": static_templates_information_html,
	"static/templates/install_step_0.html": static_templates_install_step_0_html,
	"static/templates/install_step_1.html": static_templates_install_step_1_html,
	"static/templates/install_step_2.html": static_templates_install_step_2_html,
	"static/templates/interface.html": static_templates_interface_html,
	"static/templates/list_cf_projects.tpl": static_templates_list_cf_projects_tpl,
	"static/templates/login.html": static_templates_login_html,
	"static/templates/menu.html": static_templates_menu_html,
	"static/templates/mining_menu.html": static_templates_mining_menu_html,
	"static/templates/mining_promised_amount.html": static_templates_mining_promised_amount_html,
	"static/templates/modal.html": static_templates_modal_html,
	"static/templates/money_back.html": static_templates_money_back_html,
	"static/templates/money_back_request.html": static_templates_money_back_request_html,
	"static/templates/my_cf_projects.html": static_templates_my_cf_projects_html,
	"static/templates/new_cf_project.html": static_templates_new_cf_project_html,
	"static/templates/new_credit.html": static_templates_new_credit_html,
	"static/templates/new_holidays.html": static_templates_new_holidays_html,
	"static/templates/new_promised_amount.html": static_templates_new_promised_amount_html,
	"static/templates/new_user.html": static_templates_new_user_html,
	"static/templates/node_config.html": static_templates_node_config_html,
	"static/templates/notifications.html": static_templates_notifications_html,
	"static/templates/pct.tpl": static_templates_pct_tpl,
	"static/templates/points.html": static_templates_points_html,
	"static/templates/pool_admin.html": static_templates_pool_admin_html,
	"static/templates/pool_tech_works.tpl": static_templates_pool_tech_works_tpl,
	"static/templates/progress.tpl": static_templates_progress_tpl,
	"static/templates/progress_bar.html": static_templates_progress_bar_html,
	"static/templates/promised_amount_actualization.tpl": static_templates_promised_amount_actualization_tpl,
	"static/templates/promised_amount_list.html": static_templates_promised_amount_list_html,
	"static/templates/reduction.tpl": static_templates_reduction_tpl,
	"static/templates/repayment_credit.html": static_templates_repayment_credit_html,
	"static/templates/restoring_access.html": static_templates_restoring_access_html,
	"static/templates/rewrite_primary_key.html": static_templates_rewrite_primary_key_html,
	"static/templates/sign_up_in_the_pool.html": static_templates_sign_up_in_the_pool_html,
	"static/templates/signatures.html": static_templates_signatures_html,
	"static/templates/statistic.html": static_templates_statistic_html,
	"static/templates/statistic_voting.html": static_templates_statistic_voting_html,
	"static/templates/updating_blockchain.html": static_templates_updating_blockchain_html,
	"static/templates/upgrade.html": static_templates_upgrade_html,
	"static/templates/upgrade_0.html": static_templates_upgrade_0_html,
	"static/templates/upgrade_1_and_2.html": static_templates_upgrade_1_and_2_html,
	"static/templates/upgrade_3.html": static_templates_upgrade_3_html,
	"static/templates/upgrade_4.html": static_templates_upgrade_4_html,
	"static/templates/upgrade_5.html": static_templates_upgrade_5_html,
	"static/templates/upgrade_6.html": static_templates_upgrade_6_html,
	"static/templates/upgrade_7.html": static_templates_upgrade_7_html,
	"static/templates/upgrade_resend.tpl": static_templates_upgrade_resend_tpl,
	"static/templates/vote_for_me.html": static_templates_vote_for_me_html,
	"static/templates/voting.html": static_templates_voting_html,
	"static/templates/wallets_list.html": static_templates_wallets_list_html,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"static": &_bintree_t{nil, map[string]*_bintree_t{
		"1block.bin": &_bintree_t{static_1block_bin, map[string]*_bintree_t{
		}},
		"css": &_bintree_t{nil, map[string]*_bintree_t{
			"AS3Cam.css": &_bintree_t{static_css_as3cam_css, map[string]*_bintree_t{
			}},
			"Jcrop.gif": &_bintree_t{static_css_jcrop_gif, map[string]*_bintree_t{
			}},
			"bootstrap-responsive.css": &_bintree_t{static_css_bootstrap_responsive_css, map[string]*_bintree_t{
			}},
			"bootstrap.css": &_bintree_t{static_css_bootstrap_css, map[string]*_bintree_t{
			}},
			"bootstrap.min.css": &_bintree_t{static_css_bootstrap_min_css, map[string]*_bintree_t{
			}},
			"cf.css": &_bintree_t{static_css_cf_css, map[string]*_bintree_t{
			}},
			"cropper.css": &_bintree_t{static_css_cropper_css, map[string]*_bintree_t{
			}},
			"font-awesome.css": &_bintree_t{static_css_font_awesome_css, map[string]*_bintree_t{
			}},
			"jquery-ui-timepicker-addon.css": &_bintree_t{static_css_jquery_ui_timepicker_addon_css, map[string]*_bintree_t{
			}},
			"jquery-ui.css": &_bintree_t{static_css_jquery_ui_css, map[string]*_bintree_t{
			}},
			"jquery.Jcrop.css": &_bintree_t{static_css_jquery_jcrop_css, map[string]*_bintree_t{
			}},
			"jquery.Jcrop.min.css": &_bintree_t{static_css_jquery_jcrop_min_css, map[string]*_bintree_t{
			}},
			"jquery.css": &_bintree_t{static_css_jquery_css, map[string]*_bintree_t{
			}},
			"jquery.qtip.min.css": &_bintree_t{static_css_jquery_qtip_min_css, map[string]*_bintree_t{
			}},
			"plugins": &_bintree_t{nil, map[string]*_bintree_t{
				"metisMenu": &_bintree_t{nil, map[string]*_bintree_t{
					"metisMenu.min.css": &_bintree_t{static_css_plugins_metismenu_metismenu_min_css, map[string]*_bintree_t{
					}},
				}},
			}},
			"progress.css": &_bintree_t{static_css_progress_css, map[string]*_bintree_t{
			}},
			"sb-admin-2.css": &_bintree_t{static_css_sb_admin_2_css, map[string]*_bintree_t{
			}},
			"sb-admin.css": &_bintree_t{static_css_sb_admin_css, map[string]*_bintree_t{
			}},
			"social-buttons.css": &_bintree_t{static_css_social_buttons_css, map[string]*_bintree_t{
			}},
			"tooltipster-shadow.css": &_bintree_t{static_css_tooltipster_shadow_css, map[string]*_bintree_t{
			}},
			"tooltipster.css": &_bintree_t{static_css_tooltipster_css, map[string]*_bintree_t{
			}},
		}},
		"fonts": &_bintree_t{nil, map[string]*_bintree_t{
			"FontAwesome.otf": &_bintree_t{static_fonts_fontawesome_otf, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.eot": &_bintree_t{static_fonts_fontawesome_webfont_eot, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.svg": &_bintree_t{static_fonts_fontawesome_webfont_svg, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.ttf": &_bintree_t{static_fonts_fontawesome_webfont_ttf, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.woff": &_bintree_t{static_fonts_fontawesome_webfont_woff, map[string]*_bintree_t{
			}},
			"luxisr.ttf": &_bintree_t{static_fonts_luxisr_ttf, map[string]*_bintree_t{
			}},
		}},
		"img": &_bintree_t{nil, map[string]*_bintree_t{
			"alert.png": &_bintree_t{static_img_alert_png, map[string]*_bintree_t{
			}},
			"blank.png": &_bintree_t{static_img_blank_png, map[string]*_bintree_t{
			}},
			"cf_blurb_img.png": &_bintree_t{static_img_cf_blurb_img_png, map[string]*_bintree_t{
			}},
			"face.jpg": &_bintree_t{static_img_face_jpg, map[string]*_bintree_t{
			}},
			"glyphicons-halflings.png": &_bintree_t{static_img_glyphicons_halflings_png, map[string]*_bintree_t{
			}},
			"k_bg.png": &_bintree_t{static_img_k_bg_png, map[string]*_bintree_t{
			}},
			"k_bg_pass.png": &_bintree_t{static_img_k_bg_pass_png, map[string]*_bintree_t{
			}},
			"logo-small.png": &_bintree_t{static_img_logo_small_png, map[string]*_bintree_t{
			}},
			"logo.png": &_bintree_t{static_img_logo_png, map[string]*_bintree_t{
			}},
			"noavatar.png": &_bintree_t{static_img_noavatar_png, map[string]*_bintree_t{
			}},
			"photo.png": &_bintree_t{static_img_photo_png, map[string]*_bintree_t{
			}},
			"profile.jpg": &_bintree_t{static_img_profile_jpg, map[string]*_bintree_t{
			}},
			"race.gif": &_bintree_t{static_img_race_gif, map[string]*_bintree_t{
			}},
			"us-ru.png": &_bintree_t{static_img_us_ru_png, map[string]*_bintree_t{
			}},
		}},
		"js": &_bintree_t{nil, map[string]*_bintree_t{
			"Serpent.js": &_bintree_t{static_js_serpent_js, map[string]*_bintree_t{
			}},
			"aes.js": &_bintree_t{static_js_aes_js, map[string]*_bintree_t{
			}},
			"asn1hex-1.1.js": &_bintree_t{static_js_asn1hex_1_1_js, map[string]*_bintree_t{
			}},
			"asn1hex-1.1.min.js": &_bintree_t{static_js_asn1hex_1_1_min_js, map[string]*_bintree_t{
			}},
			"asn1hex-1.js": &_bintree_t{static_js_asn1hex_1_js, map[string]*_bintree_t{
			}},
			"base64.js": &_bintree_t{static_js_base64_js, map[string]*_bintree_t{
			}},
			"bootstrap-alert.js": &_bintree_t{static_js_bootstrap_alert_js, map[string]*_bintree_t{
			}},
			"bootstrap-button.js": &_bintree_t{static_js_bootstrap_button_js, map[string]*_bintree_t{
			}},
			"bootstrap-carousel.js": &_bintree_t{static_js_bootstrap_carousel_js, map[string]*_bintree_t{
			}},
			"bootstrap-collapse.js": &_bintree_t{static_js_bootstrap_collapse_js, map[string]*_bintree_t{
			}},
			"bootstrap-dropdown.js": &_bintree_t{static_js_bootstrap_dropdown_js, map[string]*_bintree_t{
			}},
			"bootstrap-modal.js": &_bintree_t{static_js_bootstrap_modal_js, map[string]*_bintree_t{
			}},
			"bootstrap-popover.js": &_bintree_t{static_js_bootstrap_popover_js, map[string]*_bintree_t{
			}},
			"bootstrap-scrollspy.js": &_bintree_t{static_js_bootstrap_scrollspy_js, map[string]*_bintree_t{
			}},
			"bootstrap-tab.js": &_bintree_t{static_js_bootstrap_tab_js, map[string]*_bintree_t{
			}},
			"bootstrap-tooltip.js": &_bintree_t{static_js_bootstrap_tooltip_js, map[string]*_bintree_t{
			}},
			"bootstrap-transition.js": &_bintree_t{static_js_bootstrap_transition_js, map[string]*_bintree_t{
			}},
			"bootstrap-typeahead.js": &_bintree_t{static_js_bootstrap_typeahead_js, map[string]*_bintree_t{
			}},
			"bootstrap.min.js": &_bintree_t{static_js_bootstrap_min_js, map[string]*_bintree_t{
			}},
			"cropper.js": &_bintree_t{static_js_cropper_js, map[string]*_bintree_t{
			}},
			"enc-base64-min.js": &_bintree_t{static_js_enc_base64_min_js, map[string]*_bintree_t{
			}},
			"html5shiv.js": &_bintree_t{static_js_html5shiv_js, map[string]*_bintree_t{
			}},
			"index.js": &_bintree_t{static_js_index_js, map[string]*_bintree_t{
			}},
			"infobubble.js": &_bintree_t{static_js_infobubble_js, map[string]*_bintree_t{
			}},
			"jquery-1.11.0.js": &_bintree_t{static_js_jquery_1_11_0_js, map[string]*_bintree_t{
			}},
			"jquery-1.9.1.min.js": &_bintree_t{static_js_jquery_1_9_1_min_js, map[string]*_bintree_t{
			}},
			"jquery-ui-sliderAccess.js": &_bintree_t{static_js_jquery_ui_slideraccess_js, map[string]*_bintree_t{
			}},
			"jquery-ui-timepicker-addon.js": &_bintree_t{static_js_jquery_ui_timepicker_addon_js, map[string]*_bintree_t{
			}},
			"jquery-ui.min.js": &_bintree_t{static_js_jquery_ui_min_js, map[string]*_bintree_t{
			}},
			"jquery.Jcrop.js": &_bintree_t{static_js_jquery_jcrop_js, map[string]*_bintree_t{
			}},
			"jquery.js": &_bintree_t{static_js_jquery_js, map[string]*_bintree_t{
			}},
			"jquery.min.js": &_bintree_t{static_js_jquery_min_js, map[string]*_bintree_t{
			}},
			"jquery.qtip.min.js": &_bintree_t{static_js_jquery_qtip_min_js, map[string]*_bintree_t{
			}},
			"jquery.tooltipster.min.js": &_bintree_t{static_js_jquery_tooltipster_min_js, map[string]*_bintree_t{
			}},
			"jquery.webcam.as3.js": &_bintree_t{static_js_jquery_webcam_as3_js, map[string]*_bintree_t{
			}},
			"jquery_002.js": &_bintree_t{static_js_jquery_002_js, map[string]*_bintree_t{
			}},
			"js.js": &_bintree_t{static_js_js_js, map[string]*_bintree_t{
			}},
			"jsbn.js": &_bintree_t{static_js_jsbn_js, map[string]*_bintree_t{
			}},
			"jsbn2.js": &_bintree_t{static_js_jsbn2_js, map[string]*_bintree_t{
			}},
			"markerclusterer.js": &_bintree_t{static_js_markerclusterer_js, map[string]*_bintree_t{
			}},
			"mcrypt.js": &_bintree_t{static_js_mcrypt_js, map[string]*_bintree_t{
			}},
			"md5.js": &_bintree_t{static_js_md5_js, map[string]*_bintree_t{
			}},
			"plugins": &_bintree_t{nil, map[string]*_bintree_t{
				"metisMenu": &_bintree_t{nil, map[string]*_bintree_t{
					"metisMenu.min.js": &_bintree_t{static_js_plugins_metismenu_metismenu_min_js, map[string]*_bintree_t{
					}},
				}},
			}},
			"prng4.js": &_bintree_t{static_js_prng4_js, map[string]*_bintree_t{
			}},
			"rijndael.js": &_bintree_t{static_js_rijndael_js, map[string]*_bintree_t{
			}},
			"ripemd160.js": &_bintree_t{static_js_ripemd160_js, map[string]*_bintree_t{
			}},
			"rng.js": &_bintree_t{static_js_rng_js, map[string]*_bintree_t{
			}},
			"rsa.js": &_bintree_t{static_js_rsa_js, map[string]*_bintree_t{
			}},
			"rsa2.js": &_bintree_t{static_js_rsa2_js, map[string]*_bintree_t{
			}},
			"rsapem-1.1.js": &_bintree_t{static_js_rsapem_1_1_js, map[string]*_bintree_t{
			}},
			"rsapem-1.1.min.js": &_bintree_t{static_js_rsapem_1_1_min_js, map[string]*_bintree_t{
			}},
			"rsapem-1.js": &_bintree_t{static_js_rsapem_1_js, map[string]*_bintree_t{
			}},
			"rsasign-1.2.js": &_bintree_t{static_js_rsasign_1_2_js, map[string]*_bintree_t{
			}},
			"rsasign-1.2.min.js": &_bintree_t{static_js_rsasign_1_2_min_js, map[string]*_bintree_t{
			}},
			"rsasign-1.js": &_bintree_t{static_js_rsasign_1_js, map[string]*_bintree_t{
			}},
			"sAS3Cam.js": &_bintree_t{static_js_sas3cam_js, map[string]*_bintree_t{
			}},
			"sb-admin-2.js": &_bintree_t{static_js_sb_admin_2_js, map[string]*_bintree_t{
			}},
			"sb-admin.js": &_bintree_t{static_js_sb_admin_js, map[string]*_bintree_t{
			}},
			"sha1.js": &_bintree_t{static_js_sha1_js, map[string]*_bintree_t{
			}},
			"sha256.js": &_bintree_t{static_js_sha256_js, map[string]*_bintree_t{
			}},
			"sha512.js": &_bintree_t{static_js_sha512_js, map[string]*_bintree_t{
			}},
			"spin.js": &_bintree_t{static_js_spin_js, map[string]*_bintree_t{
			}},
			"spots.js": &_bintree_t{static_js_spots_js, map[string]*_bintree_t{
			}},
			"stacktable.js": &_bintree_t{static_js_stacktable_js, map[string]*_bintree_t{
			}},
			"unixtime.js": &_bintree_t{static_js_unixtime_js, map[string]*_bintree_t{
			}},
			"uploader.js": &_bintree_t{static_js_uploader_js, map[string]*_bintree_t{
			}},
			"x509-1.1.js": &_bintree_t{static_js_x509_1_1_js, map[string]*_bintree_t{
			}},
			"x509-1.1.min.js": &_bintree_t{static_js_x509_1_1_min_js, map[string]*_bintree_t{
			}},
			"youtube_webcam.js": &_bintree_t{static_js_youtube_webcam_js, map[string]*_bintree_t{
			}},
		}},
		"lang": &_bintree_t{nil, map[string]*_bintree_t{
			"1.ini": &_bintree_t{static_lang_1_ini, map[string]*_bintree_t{
			}},
			"42.ini": &_bintree_t{static_lang_42_ini, map[string]*_bintree_t{
			}},
			"en-us.all.json": &_bintree_t{static_lang_en_us_all_json, map[string]*_bintree_t{
			}},
			"locale_en-US.ini": &_bintree_t{static_lang_locale_en_us_ini, map[string]*_bintree_t{
			}},
			"locale_ru-RU.ini": &_bintree_t{static_lang_locale_ru_ru_ini, map[string]*_bintree_t{
			}},
		}},
		"nodes.inc": &_bintree_t{static_nodes_inc, map[string]*_bintree_t{
		}},
		"templates": &_bintree_t{nil, map[string]*_bintree_t{
			"abuse.html": &_bintree_t{static_templates_abuse_html, map[string]*_bintree_t{
			}},
			"add_cf_project_data.html": &_bintree_t{static_templates_add_cf_project_data_html, map[string]*_bintree_t{
			}},
			"after_install.html": &_bintree_t{static_templates_after_install_html, map[string]*_bintree_t{
			}},
			"alert_success.html": &_bintree_t{static_templates_alert_success_html, map[string]*_bintree_t{
			}},
			"arbitration.html": &_bintree_t{static_templates_arbitration_html, map[string]*_bintree_t{
			}},
			"arbitration_arbitrator.html": &_bintree_t{static_templates_arbitration_arbitrator_html, map[string]*_bintree_t{
			}},
			"arbitration_buyer.html": &_bintree_t{static_templates_arbitration_buyer_html, map[string]*_bintree_t{
			}},
			"arbitration_seller.html": &_bintree_t{static_templates_arbitration_seller_html, map[string]*_bintree_t{
			}},
			"assignments.html": &_bintree_t{static_templates_assignments_html, map[string]*_bintree_t{
			}},
			"assignments_new_miner.html": &_bintree_t{static_templates_assignments_new_miner_html, map[string]*_bintree_t{
			}},
			"assignments_promised_amount.html": &_bintree_t{static_templates_assignments_promised_amount_html, map[string]*_bintree_t{
			}},
			"block_explorer.html": &_bintree_t{static_templates_block_explorer_html, map[string]*_bintree_t{
			}},
			"bug_reporting.html": &_bintree_t{static_templates_bug_reporting_html, map[string]*_bintree_t{
			}},
			"cash_request_in.html": &_bintree_t{static_templates_cash_request_in_html, map[string]*_bintree_t{
			}},
			"cash_request_out.html": &_bintree_t{static_templates_cash_request_out_html, map[string]*_bintree_t{
			}},
			"cf_catalog.html": &_bintree_t{static_templates_cf_catalog_html, map[string]*_bintree_t{
			}},
			"cf_page_preview.html": &_bintree_t{static_templates_cf_page_preview_html, map[string]*_bintree_t{
			}},
			"cf_project_change_category.html": &_bintree_t{static_templates_cf_project_change_category_html, map[string]*_bintree_t{
			}},
			"cf_start.html": &_bintree_t{static_templates_cf_start_html, map[string]*_bintree_t{
			}},
			"change_arbitrator_conditions.html": &_bintree_t{static_templates_change_arbitrator_conditions_html, map[string]*_bintree_t{
			}},
			"change_avatar.tpl": &_bintree_t{static_templates_change_avatar_tpl, map[string]*_bintree_t{
			}},
			"change_commission.html": &_bintree_t{static_templates_change_commission_html, map[string]*_bintree_t{
			}},
			"change_country_race.tpl": &_bintree_t{static_templates_change_country_race_tpl, map[string]*_bintree_t{
			}},
			"change_creditor.html": &_bintree_t{static_templates_change_creditor_html, map[string]*_bintree_t{
			}},
			"change_geolocation.html": &_bintree_t{static_templates_change_geolocation_html, map[string]*_bintree_t{
			}},
			"change_host.html": &_bintree_t{static_templates_change_host_html, map[string]*_bintree_t{
			}},
			"change_key_close.html": &_bintree_t{static_templates_change_key_close_html, map[string]*_bintree_t{
			}},
			"change_key_request.html": &_bintree_t{static_templates_change_key_request_html, map[string]*_bintree_t{
			}},
			"change_money_back_time.html": &_bintree_t{static_templates_change_money_back_time_html, map[string]*_bintree_t{
			}},
			"change_node_key.html": &_bintree_t{static_templates_change_node_key_html, map[string]*_bintree_t{
			}},
			"change_primary_key.html": &_bintree_t{static_templates_change_primary_key_html, map[string]*_bintree_t{
			}},
			"change_promised_amount.html": &_bintree_t{static_templates_change_promised_amount_html, map[string]*_bintree_t{
			}},
			"credits.html": &_bintree_t{static_templates_credits_html, map[string]*_bintree_t{
			}},
			"currency_exchange.html": &_bintree_t{static_templates_currency_exchange_html, map[string]*_bintree_t{
			}},
			"currency_exchange_delete.html": &_bintree_t{static_templates_currency_exchange_delete_html, map[string]*_bintree_t{
			}},
			"db_info.html": &_bintree_t{static_templates_db_info_html, map[string]*_bintree_t{
			}},
			"del_cf_funding.html": &_bintree_t{static_templates_del_cf_funding_html, map[string]*_bintree_t{
			}},
			"del_cf_project.html": &_bintree_t{static_templates_del_cf_project_html, map[string]*_bintree_t{
			}},
			"del_credit.html": &_bintree_t{static_templates_del_credit_html, map[string]*_bintree_t{
			}},
			"del_promised_amount.html": &_bintree_t{static_templates_del_promised_amount_html, map[string]*_bintree_t{
			}},
			"for_repaid_fix.html": &_bintree_t{static_templates_for_repaid_fix_html, map[string]*_bintree_t{
			}},
			"holidays_list.html": &_bintree_t{static_templates_holidays_list_html, map[string]*_bintree_t{
			}},
			"home.html": &_bintree_t{static_templates_home_html, map[string]*_bintree_t{
			}},
			"index.html": &_bintree_t{static_templates_index_html, map[string]*_bintree_t{
			}},
			"index_cf.html": &_bintree_t{static_templates_index_cf_html, map[string]*_bintree_t{
			}},
			"information.html": &_bintree_t{static_templates_information_html, map[string]*_bintree_t{
			}},
			"install_step_0.html": &_bintree_t{static_templates_install_step_0_html, map[string]*_bintree_t{
			}},
			"install_step_1.html": &_bintree_t{static_templates_install_step_1_html, map[string]*_bintree_t{
			}},
			"install_step_2.html": &_bintree_t{static_templates_install_step_2_html, map[string]*_bintree_t{
			}},
			"interface.html": &_bintree_t{static_templates_interface_html, map[string]*_bintree_t{
			}},
			"list_cf_projects.tpl": &_bintree_t{static_templates_list_cf_projects_tpl, map[string]*_bintree_t{
			}},
			"login.html": &_bintree_t{static_templates_login_html, map[string]*_bintree_t{
			}},
			"menu.html": &_bintree_t{static_templates_menu_html, map[string]*_bintree_t{
			}},
			"mining_menu.html": &_bintree_t{static_templates_mining_menu_html, map[string]*_bintree_t{
			}},
			"mining_promised_amount.html": &_bintree_t{static_templates_mining_promised_amount_html, map[string]*_bintree_t{
			}},
			"modal.html": &_bintree_t{static_templates_modal_html, map[string]*_bintree_t{
			}},
			"money_back.html": &_bintree_t{static_templates_money_back_html, map[string]*_bintree_t{
			}},
			"money_back_request.html": &_bintree_t{static_templates_money_back_request_html, map[string]*_bintree_t{
			}},
			"my_cf_projects.html": &_bintree_t{static_templates_my_cf_projects_html, map[string]*_bintree_t{
			}},
			"new_cf_project.html": &_bintree_t{static_templates_new_cf_project_html, map[string]*_bintree_t{
			}},
			"new_credit.html": &_bintree_t{static_templates_new_credit_html, map[string]*_bintree_t{
			}},
			"new_holidays.html": &_bintree_t{static_templates_new_holidays_html, map[string]*_bintree_t{
			}},
			"new_promised_amount.html": &_bintree_t{static_templates_new_promised_amount_html, map[string]*_bintree_t{
			}},
			"new_user.html": &_bintree_t{static_templates_new_user_html, map[string]*_bintree_t{
			}},
			"node_config.html": &_bintree_t{static_templates_node_config_html, map[string]*_bintree_t{
			}},
			"notifications.html": &_bintree_t{static_templates_notifications_html, map[string]*_bintree_t{
			}},
			"pct.tpl": &_bintree_t{static_templates_pct_tpl, map[string]*_bintree_t{
			}},
			"points.html": &_bintree_t{static_templates_points_html, map[string]*_bintree_t{
			}},
			"pool_admin.html": &_bintree_t{static_templates_pool_admin_html, map[string]*_bintree_t{
			}},
			"pool_tech_works.tpl": &_bintree_t{static_templates_pool_tech_works_tpl, map[string]*_bintree_t{
			}},
			"progress.tpl": &_bintree_t{static_templates_progress_tpl, map[string]*_bintree_t{
			}},
			"progress_bar.html": &_bintree_t{static_templates_progress_bar_html, map[string]*_bintree_t{
			}},
			"promised_amount_actualization.tpl": &_bintree_t{static_templates_promised_amount_actualization_tpl, map[string]*_bintree_t{
			}},
			"promised_amount_list.html": &_bintree_t{static_templates_promised_amount_list_html, map[string]*_bintree_t{
			}},
			"reduction.tpl": &_bintree_t{static_templates_reduction_tpl, map[string]*_bintree_t{
			}},
			"repayment_credit.html": &_bintree_t{static_templates_repayment_credit_html, map[string]*_bintree_t{
			}},
			"restoring_access.html": &_bintree_t{static_templates_restoring_access_html, map[string]*_bintree_t{
			}},
			"rewrite_primary_key.html": &_bintree_t{static_templates_rewrite_primary_key_html, map[string]*_bintree_t{
			}},
			"sign_up_in_the_pool.html": &_bintree_t{static_templates_sign_up_in_the_pool_html, map[string]*_bintree_t{
			}},
			"signatures.html": &_bintree_t{static_templates_signatures_html, map[string]*_bintree_t{
			}},
			"statistic.html": &_bintree_t{static_templates_statistic_html, map[string]*_bintree_t{
			}},
			"statistic_voting.html": &_bintree_t{static_templates_statistic_voting_html, map[string]*_bintree_t{
			}},
			"updating_blockchain.html": &_bintree_t{static_templates_updating_blockchain_html, map[string]*_bintree_t{
			}},
			"upgrade.html": &_bintree_t{static_templates_upgrade_html, map[string]*_bintree_t{
			}},
			"upgrade_0.html": &_bintree_t{static_templates_upgrade_0_html, map[string]*_bintree_t{
			}},
			"upgrade_1_and_2.html": &_bintree_t{static_templates_upgrade_1_and_2_html, map[string]*_bintree_t{
			}},
			"upgrade_3.html": &_bintree_t{static_templates_upgrade_3_html, map[string]*_bintree_t{
			}},
			"upgrade_4.html": &_bintree_t{static_templates_upgrade_4_html, map[string]*_bintree_t{
			}},
			"upgrade_5.html": &_bintree_t{static_templates_upgrade_5_html, map[string]*_bintree_t{
			}},
			"upgrade_6.html": &_bintree_t{static_templates_upgrade_6_html, map[string]*_bintree_t{
			}},
			"upgrade_7.html": &_bintree_t{static_templates_upgrade_7_html, map[string]*_bintree_t{
			}},
			"upgrade_resend.tpl": &_bintree_t{static_templates_upgrade_resend_tpl, map[string]*_bintree_t{
			}},
			"vote_for_me.html": &_bintree_t{static_templates_vote_for_me_html, map[string]*_bintree_t{
			}},
			"voting.html": &_bintree_t{static_templates_voting_html, map[string]*_bintree_t{
			}},
			"wallets_list.html": &_bintree_t{static_templates_wallets_list_html, map[string]*_bintree_t{
			}},
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

