---

name: omnilog-ui-designer
description: Use this when implementing or improving UI in Omnilog, fixing CSS/layout issues, creating new pages, or making the interface more consistent with the existing design system.

---

# Purpose

Use this skill when:

* creating new pages or components for Omnilog;
* improving existing UI without a provided mockup;
* fixing CSS bugs, layout inconsistencies, or responsive issues;
* refactoring existing styles for better maintainability;
* improving UX while preserving the current Omnilog visual identity.

# Omnilog design philosophy

The interface should look like a professional SaaS logistics platform.

Design priorities:

* clean;
* minimalistic;
* modern;
* consistent;
* data-focused;
* business-oriented.

Avoid decorative elements that do not improve usability.

Every page should feel like part of the same application.

# Existing visual language

Always preserve the existing Omnilog style:

* dark theme;
* blue accent colors;
* rounded cards;
* subtle borders;
* soft shadows;
* large readable headings;
* generous spacing between sections;
* consistent paddings and margins;
* uniform button heights;
* uniform table styling;
* consistent badges and status labels.

Do not introduce another design language.

# Reuse rules

Before creating anything new:

1. Search existing components.
2. Search existing layouts.
3. Search existing CSS utilities.
4. Search shared variables and tokens.

Reuse existing components whenever possible.

Never duplicate an existing component because it is "almost suitable".

Refactor existing components instead.

# CSS rules

Prefer existing project CSS architecture.

Reuse:

* CSS variables;
* spacing tokens;
* border radius values;
* shadows;
* typography;
* breakpoints;
* animations;
* colors.

Never hardcode values that already exist globally.

If a new design token is required, add it globally instead of locally.

# Layout rules

All pages should follow the same spacing system.

Maintain consistent:

* top padding;
* left/right content padding;
* card spacing;
* table spacing;
* section spacing;
* page width.

Avoid pages looking visually shifted compared to others.

# Scrollbar rules

Scrollbar should exist **only when content genuinely overflows**.

Never create scrollbars because of:

* incorrect `100vh`;
* incorrect `100dvh`;
* extra margins;
* extra padding;
* transforms;
* hidden wrappers;
* loading placeholders;
* nested overflow containers.

Avoid nested scrolling whenever possible.

The page should scroll naturally.

Horizontal scrollbars should never appear unless the content itself requires them.

# Responsive rules

Desktop-first but responsive.

Support:

* 1920px
* 1600px
* 1440px
* 1366px
* 1280px
* 1024px

Cards and grids should adapt smoothly without breaking spacing.

Avoid sudden wrapping that changes layout dramatically.

# Tables

Tables should:

* have aligned columns;
* have consistent row height;
* keep action buttons aligned;
* avoid unnecessary wrapping;
* preserve visual rhythm.

Action buttons should preferably use:

* Support button
* Overflow menu (⋮)

instead of multiple action buttons occupying horizontal space.

# Sidebar rules

Sidebar items must:

* keep fixed height;
* never jump to two lines on hover;
* never resize during hover;
* keep icon and text aligned;
* preserve width during transitions.

Hover effects should never change layout.

# Navigation transitions

All pages should use the same transition system.

Avoid:

* flashing;
* blinking;
* layout jumps;
* scrollbar appearing/disappearing;
* content shifting sideways.

Navigation should feel smooth and consistent.

# Loading behavior

Loading placeholders must preserve layout dimensions.

During loading:

* no layout shift;
* no scrollbar flickering;
* no width changes;
* no content jumping.

Initial render should match final render dimensions.

# Visual consistency

Cards should:

* share identical border radius;
* share shadows;
* share border opacity;
* share spacing;
* share typography hierarchy.

Buttons should:

* have equal heights;
* consistent padding;
* consistent hover effects;
* consistent disabled states.

Forms should:

* have identical input heights;
* identical label spacing;
* identical validation messages.

# Implementation approach

1. Inspect existing components.
2. Inspect existing layout wrappers.
3. Inspect global styles.
4. Reuse existing design tokens.
5. Fix the root cause instead of patching symptoms.
6. Preserve visual consistency across the whole application.
7. Avoid hacks, magic numbers, and page-specific CSS unless absolutely necessary.

# Expected output

The resulting UI should look like a polished commercial logistics SaaS platform.

Every page should feel visually unified, with consistent spacing, typography, cards, tables, forms, navigation, transitions, and responsive behavior.

No unnecessary scrollbars, layout shifts, flickering, or inconsistent styling should remain.
