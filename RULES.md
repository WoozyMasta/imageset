<!-- Automatically generated file, do not modify! -->

# Lint Rules Registry

This document contains the current registry of lint rules.

Total rules: 22.

## imageset

ImageSet

> Lint rules for image set atlas and entry validation.

### validate

> Image set semantic validation diagnostics.

Codes:
[IMGSET2001](#imgset2001),
[IMGSET2002](#imgset2002),
[IMGSET2003](#imgset2003),
[IMGSET2004](#imgset2004),
[IMGSET2005](#imgset2005),
[IMGSET2006](#imgset2006),
[IMGSET2007](#imgset2007),
[IMGSET2008](#imgset2008),
[IMGSET2009](#imgset2009),
[IMGSET2010](#imgset2010),
[IMGSET2011](#imgset2011),
[IMGSET2012](#imgset2012),
[IMGSET2013](#imgset2013),
[IMGSET2014](#imgset2014),
[IMGSET2015](#imgset2015),
[IMGSET2016](#imgset2016),
[IMGSET2017](#imgset2017),
[IMGSET2018](#imgset2018),
[IMGSET2019](#imgset2019),
[IMGSET2020](#imgset2020),
[IMGSET2021](#imgset2021),
[IMGSET2022](#imgset2022),

#### `IMGSET2001`

`Ref_size.width` must be greater than 0

> Root atlas width is required and must be a positive integer.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.ref-size-width-must-be-greater-than-0` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2002`

`Ref_size.height` must be greater than 0

> Root atlas height is required and must be a positive integer.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.ref-size-height-must-be-greater-than-0` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2003`

`Ref_size` side should be a power of two

> Power-of-two size improves compatibility with common atlas workflows.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.ref-size-side-should-be-a-power-of-two` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `IMGSET2004`

`Textures` section must contain at least one texture

> DayZ expects texture list to exist and contain at least one entry.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.textures-section-must-contain-at-least-one-texture` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2005`

`Texture.path` must be non-empty

> Each texture entry must define a path to texture source file.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.texture-path-must-be-non-empty` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2006`

`Texture.mpix` must be zero or greater

> Negative pixels-per-meter value is invalid for texture metadata.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.texture-mpix-must-be-zero-or-greater` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2007`

`Images` section is empty

> Currently warning-only to allow group-only files during migration.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.images-section-is-empty` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `IMGSET2008`

`Group` name must be non-empty

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.group-name-must-be-non-empty` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2009`

`Group` name must be unique

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.group-name-must-be-unique` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2010`

`Group.images` section is empty

> Group exists but does not define any image entries.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.group-images-section-is-empty` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `IMGSET2011`

`Image` name must be non-empty

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-name-must-be-non-empty` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2012`

`Image` name must be unique in current section

> Duplicates are checked separately for root images and for each group.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-name-must-be-unique-in-current-section` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2013`

`Image` name is duplicated globally across root/groups

> Informational only: global duplicate may still resolve by group prefix.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-name-is-duplicated-globally-across-root-groups` |
| Scope | `validate` |
| Severity | `info` |
| Enabled | `true` (implicit) |

#### `IMGSET2014`

`Image.pos.x` must be zero or greater

> Negative coordinates place sprite origin outside atlas bounds.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-pos-x-must-be-zero-or-greater` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2015`

`Image.pos.y` must be zero or greater

> Negative coordinates place sprite origin outside atlas bounds.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-pos-y-must-be-zero-or-greater` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2016`

`Image.size.width` must be greater than 0

> Sprite width must be positive for usable atlas entry.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-size-width-must-be-greater-than-0` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2017`

`Image.size.height` must be greater than 0

> Sprite height must be positive for usable atlas entry.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-size-height-must-be-greater-than-0` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2018`

Image exceeds atlas bounds on width

> Triggered when pos.x + size.width is greater than ref_size.width.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-exceeds-atlas-bounds-on-width` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2019`

Image exceeds atlas bounds on height

> Triggered when pos.y + size.height is greater than ref_size.height.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-exceeds-atlas-bounds-on-height` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2020`

`Image.flags` must use supported mask values (0..3)

> Allowed flags are 0, ISHorizontalTile(1), ISVerticalTile(2), and 3.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-flags-must-use-supported-mask-values-0-3` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `IMGSET2021`

Image rectangle overlaps another image

> Overlap means two sprite rectangles intersect in atlas space.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-rectangle-overlaps-another-image` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `IMGSET2022`

Image padding is smaller than configured minimum

> Disabled by default. Enable to enforce minimum gap between image rectangles.

| Field | Value |
| --- | --- |
| Rule ID | `imageset.validate.image-padding-is-smaller-than-configured-minimum` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `false` |

Default options:
```json
{
  "min_padding": 4
}
```

---

> Generated with
> [lintkit](https://github.com/woozymasta/lintkit)
> version `dev`
> commit `unknown`

<!-- Automatically generated file, do not modify! -->
