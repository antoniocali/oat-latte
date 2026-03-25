import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
    docsSidebar: [
        'intro',
        'installation',
        {
            type: 'category',
            label: 'Tutorial',
            link: {type: 'doc', id: 'tutorial/index'},
            items: [],
        },
        'concepts',
        {
            type: 'category',
            label: 'Layout',
            link: {type: 'doc', id: 'layout'},
            items: [],
        },
        {
            type: 'category',
            label: 'Widgets',
            items: [
                'widgets/text',
                'widgets/button',
                'widgets/edittext',
                'widgets/list',
                'widgets/checkbox',
                'widgets/label',
                'widgets/progressbar',
                'widgets/notifications',
                'widgets/statusbar',
                'widgets/dialog'
            ],
        },
        'themes',
        'focus',
        'custom-components',
        'agents',
        'changelog'
    ],
};

export default sidebars;
